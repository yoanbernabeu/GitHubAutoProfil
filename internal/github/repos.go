package github

import (
	"fmt"
	"time"

	gh "github.com/google/go-github/v68/github"
)

// Repo holds relevant repository data.
type Repo struct {
	Name        string
	Description string
	URL         string
	Stars       int
	Forks       int
	Language    string
	Languages   map[string]int
	PushedAt    time.Time
	Fork        bool
	Archived    bool
}

// FetchAllRepos retrieves all owned, non-fork repositories with pagination.
func (c *Client) FetchAllRepos() ([]Repo, error) {
	var allRepos []Repo
	opts := &gh.RepositoryListByUserOptions{
		Type: "owner",
		Sort: "pushed",
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := c.API.Repositories.ListByUser(c.Context(), c.Username, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repos: %w", err)
		}

		for _, r := range repos {
			if r.GetFork() || r.GetArchived() {
				continue
			}
			repo := Repo{
				Name:        r.GetName(),
				Description: r.GetDescription(),
				URL:         r.GetHTMLURL(),
				Stars:       r.GetStargazersCount(),
				Forks:       r.GetForksCount(),
				Language:    r.GetLanguage(),
				Fork:        r.GetFork(),
				Archived:    r.GetArchived(),
			}
			if r.PushedAt != nil {
				repo.PushedAt = r.PushedAt.Time
			}
			allRepos = append(allRepos, repo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// FetchRepoLanguages retrieves the language breakdown for a repository.
func (c *Client) FetchRepoLanguages(repoName string) (map[string]int, error) {
	langs, _, err := c.API.Repositories.ListLanguages(c.Context(), c.Username, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages for %s: %w", repoName, err)
	}
	return langs, nil
}

// EnsureProfileRepo checks if the username/username repo exists, creates it if not.
func (c *Client) EnsureProfileRepo() error {
	_, _, err := c.API.Repositories.Get(c.Context(), c.Username, c.Username)
	if err != nil {
		// Repo doesn't exist, create it
		repo := &gh.Repository{
			Name:        gh.Ptr(c.Username),
			Description: gh.Ptr("My GitHub profile"),
			AutoInit:    gh.Ptr(true),
		}
		_, _, createErr := c.API.Repositories.Create(c.Context(), "", repo)
		if createErr != nil {
			return fmt.Errorf("failed to create profile repo: %w", createErr)
		}
	}
	return nil
}

// PushReadme creates or updates the README.md in the profile repo.
func (c *Client) PushReadme(content string) error {
	ctx := c.Context()
	path := "README.md"

	// Check if file already exists to get its SHA
	var sha *string
	existing, _, _, err := c.API.Repositories.GetContents(ctx, c.Username, c.Username, path, nil)
	if err == nil && existing != nil {
		sha = existing.SHA
	}

	opts := &gh.RepositoryContentFileOptions{
		Message: gh.Ptr("Update profile README"),
		Content: []byte(content),
		SHA:     sha,
	}

	_, _, err = c.API.Repositories.CreateFile(ctx, c.Username, c.Username, path, opts)
	if err != nil {
		return fmt.Errorf("failed to push README: %w", err)
	}

	return nil
}
