package analyzer

import (
	"sort"

	gh "github.com/yoanbernabeu/GitHubAutoProfil/internal/github"
)

// Stats holds aggregate statistics across all repos.
type Stats struct {
	TotalRepos    int
	TotalStars    int
	TotalForks    int
}

// ComputeStats calculates aggregate stats from all repos.
func ComputeStats(repos []gh.Repo) Stats {
	s := Stats{TotalRepos: len(repos)}
	for _, r := range repos {
		s.TotalStars += r.Stars
		s.TotalForks += r.Forks
	}
	return s
}

// TopByStars returns the top N repos sorted by star count.
func TopByStars(repos []gh.Repo, n int) []gh.Repo {
	sorted := make([]gh.Repo, len(repos))
	copy(sorted, repos)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Stars > sorted[j].Stars
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// RecentByPush returns the top N repos sorted by most recent push.
func RecentByPush(repos []gh.Repo, n int) []gh.Repo {
	sorted := make([]gh.Repo, len(repos))
	copy(sorted, repos)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PushedAt.After(sorted[j].PushedAt)
	})
	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}
