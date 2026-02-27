package github

import (
	"context"

	gh "github.com/google/go-github/v68/github"
)

// Client wraps the GitHub API client with the target username.
type Client struct {
	API      *gh.Client
	Username string
}

// NewClient creates a GitHub client for the given username.
// If token is empty, the client is unauthenticated (public API only, lower rate limits).
func NewClient(token, username string) *Client {
	api := gh.NewClient(nil)
	if token != "" {
		api = api.WithAuthToken(token)
	}
	return &Client{
		API:      api,
		Username: username,
	}
}

// Context returns a background context.
func (c *Client) Context() context.Context {
	return context.Background()
}
