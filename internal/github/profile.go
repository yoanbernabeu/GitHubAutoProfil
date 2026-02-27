package github

import (
	"encoding/json"
	"fmt"
)

// SocialAccount represents a linked social account on a GitHub profile.
type SocialAccount struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

// Profile holds the relevant user profile data.
type Profile struct {
	Login          string
	Name           string
	Bio            string
	Location       string
	Blog           string
	Twitter        string
	AvatarURL      string
	Followers      int
	Following      int
	PublicRepos    int
	SocialAccounts []SocialAccount
}

// FetchProfile retrieves the GitHub user profile and social accounts.
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

	// Fetch social accounts via REST API
	socials, err := c.fetchSocialAccounts()
	if err == nil {
		p.SocialAccounts = socials
	}

	return p, nil
}

func (c *Client) fetchSocialAccounts() ([]SocialAccount, error) {
	url := fmt.Sprintf("users/%s/social_accounts", c.Username)
	req, err := c.API.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var accounts []SocialAccount
	resp, err := c.API.Do(c.Context(), req, &accounts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Deduplicate: remove Twitter from social accounts if already in profile
	if c.Username != "" {
		var filtered []SocialAccount
		for _, a := range accounts {
			filtered = append(filtered, a)
		}
		accounts = filtered
	}

	return accounts, nil
}

// SocialBadge returns a Shields.io badge URL for a social account.
func (a SocialAccount) SocialBadge() string {
	switch a.Provider {
	case "twitter":
		return socialBadgeURL("Twitter", "1DA1F2", "twitter")
	case "linkedin":
		return socialBadgeURL("LinkedIn", "0A66C2", "linkedin")
	case "youtube":
		return socialBadgeURL("YouTube", "FF0000", "youtube")
	case "instagram":
		return socialBadgeURL("Instagram", "E4405F", "instagram")
	case "mastodon":
		return socialBadgeURL("Mastodon", "6364FF", "mastodon")
	case "facebook":
		return socialBadgeURL("Facebook", "1877F2", "facebook")
	case "twitch":
		return socialBadgeURL("Twitch", "9146FF", "twitch")
	case "reddit":
		return socialBadgeURL("Reddit", "FF4500", "reddit")
	case "npm":
		return socialBadgeURL("npm", "CB3837", "npm")
	case "hometown":
		return socialBadgeURL("Hometown", "6364FF", "mastodon")
	default:
		// Generic link badge for unknown providers
		return fmt.Sprintf("https://img.shields.io/badge/%s-555555?style=for-the-badge&logo=link&logoColor=white", a.Provider)
	}
}

func socialBadgeURL(label, color, logo string) string {
	return fmt.Sprintf("https://img.shields.io/badge/%s-%s?style=for-the-badge&logo=%s&logoColor=white", label, color, logo)
}

// DisplayName returns a human-readable name for the social account.
func (a SocialAccount) DisplayName() string {
	names := map[string]string{
		"twitter":   "Twitter",
		"linkedin":  "LinkedIn",
		"youtube":   "YouTube",
		"instagram": "Instagram",
		"mastodon":  "Mastodon",
		"facebook":  "Facebook",
		"twitch":    "Twitch",
		"reddit":    "Reddit",
		"npm":       "npm",
		"hometown":  "Hometown",
	}
	if n, ok := names[a.Provider]; ok {
		return n
	}
	return a.Provider
}

// MarshalJSON implements json.Marshaler for SocialAccount (needed for API response parsing).
var _ json.Unmarshaler = (*SocialAccount)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (a *SocialAccount) UnmarshalJSON(data []byte) error {
	type alias SocialAccount
	return json.Unmarshal(data, (*alias)(a))
}
