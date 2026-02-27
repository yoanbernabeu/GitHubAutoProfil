package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup for GitHubAutoProfil",
	Long:  "Guides you through setting up your GitHub token and username.",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to GitHubAutoProfil!")
	fmt.Println()

	// Username
	fmt.Print("Enter your GitHub username: ")
	u, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	u = strings.TrimSpace(u)
	if u == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Token
	fmt.Println()
	fmt.Println("You need a GitHub Personal Access Token with 'repo' scope.")
	fmt.Println("Create one at: https://github.com/settings/tokens")
	fmt.Print("Enter your GitHub token: ")
	t, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	t = strings.TrimSpace(t)
	if t == "" {
		return fmt.Errorf("token cannot be empty")
	}

	fmt.Println()
	fmt.Println("Configuration complete! You can now run:")
	fmt.Println()
	fmt.Printf("  export GITHUB_TOKEN=\"%s\"\n", t)
	fmt.Printf("  githubautoprofil generate --username %s --preview\n", u)
	fmt.Println()
	fmt.Println("Or to generate and push directly:")
	fmt.Printf("  githubautoprofil generate --username %s --push\n", u)
	fmt.Println()
	fmt.Println("Tip: Add the export line to your shell profile (~/.bashrc, ~/.zshrc) to persist it.")

	return nil
}
