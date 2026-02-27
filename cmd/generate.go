package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/GitHubAutoProfil/internal/analyzer"
	gh "github.com/yoanbernabeu/GitHubAutoProfil/internal/github"
	"github.com/yoanbernabeu/GitHubAutoProfil/internal/renderer"
)

var (
	preview bool
	push    bool
	top     int
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate your GitHub profile README",
	Long:  "Fetches your GitHub profile and repositories, analyzes them, and generates a README.md.",
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().BoolVar(&preview, "preview", false, "Preview the generated README locally (no push)")
	generateCmd.Flags().BoolVar(&push, "push", false, "Push the generated README to your profile repo")
	generateCmd.Flags().IntVar(&top, "top", 5, "Number of top projects to display")
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	tkn := resolveToken()
	if username == "" {
		return fmt.Errorf("GitHub username required: use --username flag")
	}
	if push && tkn == "" {
		return fmt.Errorf("GitHub token required for --push: use --token flag or set GITHUB_TOKEN")
	}

	client := gh.NewClient(tkn, username)

	// Fetch profile
	fmt.Println("Fetching profile...")
	profile, err := client.FetchProfile()
	if err != nil {
		return err
	}

	// Fetch repos
	fmt.Println("Fetching repositories...")
	repos, err := client.FetchAllRepos()
	if err != nil {
		return err
	}
	fmt.Printf("Found %d repositories\n", len(repos))

	// Fetch languages for each repo
	fmt.Println("Fetching languages...")
	var langMaps []map[string]int
	for _, r := range repos {
		langs, err := client.FetchRepoLanguages(r.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch languages for %s: %v\n", r.Name, err)
			continue
		}
		langMaps = append(langMaps, langs)
	}

	// Analyze
	stats := analyzer.ComputeStats(repos)
	topRepos := analyzer.TopByStars(repos, top)
	recentRepos := analyzer.RecentByPush(repos, 5)
	languages := analyzer.AggregateLanguages(langMaps)

	// Render
	fmt.Println("Rendering README...")
	data := renderer.Data{
		Profile:     profile,
		Stats:       stats,
		TopRepos:    topRepos,
		RecentRepos: recentRepos,
		Languages:   languages,
	}

	output, err := renderer.Render(data)
	if err != nil {
		return err
	}

	if preview || !push {
		fmt.Println("\n--- Generated README ---")
		fmt.Println(output)
		fmt.Println("--- End of README ---")
		return nil
	}

	// Push
	fmt.Println("Ensuring profile repo exists...")
	if err := client.EnsureProfileRepo(); err != nil {
		return err
	}

	fmt.Println("Pushing README...")
	if err := client.PushReadme(output); err != nil {
		return err
	}

	fmt.Printf("README successfully pushed to https://github.com/%s/%s\n", username, username)
	return nil
}
