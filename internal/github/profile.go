package github

import "fmt"

// Profile holds the relevant user profile data.
type Profile struct {
	Login     string
	Name      string
	Bio       string
	Location  string
	Blog      string
	Twitter   string
	AvatarURL string
	Followers int
	Following int
	PublicRepos int
}

// FetchProfile retrieves the GitHub user profile.
func (c *Client) FetchProfile() (*Profile, error) {
	user, _, err := c.API.Users.Get(c.Context(), c.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile for %s: %w", c.Username, err)
	}

	p := &Profile{
		Login:       user.GetLogin(),
		Name:        user.GetName(),
		Bio:         user.GetBio(),
		Location:    user.GetLocation(),
		Blog:        user.GetBlog(),
		Twitter:     user.GetTwitterUsername(),
		AvatarURL:   user.GetAvatarURL(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
		PublicRepos: user.GetPublicRepos(),
	}

	return p, nil
}
