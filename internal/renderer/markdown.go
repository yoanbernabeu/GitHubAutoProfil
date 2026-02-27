package renderer

import (
	"bytes"
	"embed"
	"fmt"
	"net/url"
	"text/template"
	"time"

	"github.com/yoanbernabeu/GitHubAutoProfil/internal/analyzer"
	gh "github.com/yoanbernabeu/GitHubAutoProfil/internal/github"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Data holds all data passed to the template.
type Data struct {
	Profile     *gh.Profile
	Stats       analyzer.Stats
	TopRepos    []gh.Repo
	RecentRepos []gh.Repo
	Languages   []analyzer.Language
}

// Render executes the default template with the provided data.
func Render(data Data) (string, error) {
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"now": func() string {
			return time.Now().UTC().Format("Jan 2, 2006")
		},
		"truncate": func(s string, n int) string {
			if len(s) <= n {
				return s
			}
			return s[:n] + "..."
		},
		"inc": func(i int) int {
			return i + 1
		},
		"urlPathEscape": func(s string) string {
			return url.PathEscape(s)
		},
	}

	tmpl, err := template.New("default.md.tmpl").Funcs(funcMap).ParseFS(templateFS, "templates/default.md.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
