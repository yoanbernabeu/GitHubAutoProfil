package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	token    string
	username string
)

var rootCmd = &cobra.Command{
	Use:   "githubautoprofil",
	Short: "Auto-generate your GitHub profile README",
	Long:  "GitHubAutoProfil scans your GitHub profile and generates a beautiful README.md for your username/username repository.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "GitHub personal access token (or set GITHUB_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "GitHub username")
}

func resolveToken() string {
	if token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}
