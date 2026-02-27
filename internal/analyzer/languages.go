package analyzer

import (
	"fmt"
	"net/url"
	"sort"
)

// Language represents a language with its usage percentage.
type Language struct {
	Name       string
	RepoCount  int
	Percentage float64
	BadgeURL   string
	Color      string
}

// AggregateLanguages normalizes language data per-repo then averages across repos.
// Each repo has equal weight regardless of size, so a large C project
// won't dominate over many smaller PHP/Go projects.
func AggregateLanguages(langMaps []map[string]int) []Language {
	if len(langMaps) == 0 {
		return nil
	}

	// For each repo, compute per-repo percentages, then sum them
	scores := make(map[string]float64)
	repoCount := make(map[string]int)

	for _, m := range langMaps {
		var repoTotal int
		for _, bytes := range m {
			repoTotal += bytes
		}
		if repoTotal == 0 {
			continue
		}
		for lang, bytes := range m {
			scores[lang] += float64(bytes) / float64(repoTotal)
			repoCount[lang]++
		}
	}

	// Normalize so all percentages sum to 100
	var totalScore float64
	for _, s := range scores {
		totalScore += s
	}
	if totalScore == 0 {
		return nil
	}

	var langs []Language
	for name, score := range scores {
		pct := score / totalScore * 100
		if pct < 1 {
			continue
		}
		color := languageColor(name)
		langs = append(langs, Language{
			Name:       name,
			RepoCount:  repoCount[name],
			Percentage: pct,
			BadgeURL:   buildBadgeURL(name, pct, color),
			Color:      color,
		})
	}

	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Percentage > langs[j].Percentage
	})

	return langs
}

func buildBadgeURL(name string, pct float64, color string) string {
	label := url.PathEscape(name)
	message := url.PathEscape(fmt.Sprintf("%.1f%%", pct))
	return fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=flat-square", label, message, color)
}

func languageColor(name string) string {
	colors := map[string]string{
		"Go":         "00ADD8",
		"JavaScript": "F7DF1E",
		"TypeScript": "3178C6",
		"Python":     "3776AB",
		"Rust":       "DEA584",
		"Java":       "ED8B00",
		"C":          "A8B9CC",
		"C++":        "00599C",
		"C#":         "239120",
		"Ruby":       "CC342D",
		"PHP":        "777BB4",
		"Swift":      "F05138",
		"Kotlin":     "7F52FF",
		"Dart":       "0175C2",
		"Shell":      "89E051",
		"HTML":       "E34F26",
		"CSS":        "1572B6",
		"Sass":       "CC6699",
		"Vue":        "4FC08D",
		"Svelte":     "FF3E00",
		"Lua":        "000080",
		"Elixir":     "6E4A7E",
		"Haskell":    "5D4F85",
		"Scala":      "DC322F",
		"R":          "276DC3",
		"MATLAB":     "E16737",
		"Jupyter Notebook": "F37626",
		"Dockerfile": "384D54",
		"Makefile":   "427819",
		"HCL":        "844FBA",
		"Nix":        "7E7EFF",
		"Zig":        "F7A41D",
	}
	if c, ok := colors[name]; ok {
		return c
	}
	return "555555"
}
