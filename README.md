<div align="center">

# GitHubAutoProfil

**Auto-generate a beautiful GitHub profile README from your repos and activity.**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/yoanbernabeu/GitHubAutoProfil?style=flat-square)](https://github.com/yoanbernabeu/GitHubAutoProfil/releases)

Scan your GitHub profile, analyze your repositories, and generate a polished `README.md` for your `username/username` repo — the one that shows up on your GitHub profile page.

[Installation](#installation) · [Usage](#usage) · [GitHub Action](#github-action) · [How it works](#how-it-works)

</div>

---

## Features

- **Profile & stats** — name, bio, location, links, followers, stars, forks
- **Top projects** — ranked by stars with descriptions
- **Recent activity** — latest pushed repos
- **Language badges** — Shields.io badges with fair percentages (normalized per-repo so one large project doesn't skew everything)
- **Auto-push** — commit the generated README directly to your profile repo
- **GitHub Action** — schedule weekly updates with zero maintenance

## Installation

### From releases (recommended)

Download the latest binary from the [releases page](https://github.com/yoanbernabeu/GitHubAutoProfil/releases).

### From source

```bash
go install github.com/yoanbernabeu/GitHubAutoProfil@latest
```

### Build locally

```bash
git clone https://github.com/yoanbernabeu/GitHubAutoProfil.git
cd GitHubAutoProfil
go build -o githubautoprofil .
```

## Usage

### Preview (no write, no token required)

```bash
githubautoprofil generate --username YOUR_USERNAME --preview
```

### Generate and push to GitHub

```bash
export GITHUB_TOKEN="ghp_xxx"
githubautoprofil generate --username YOUR_USERNAME --push
```

### Options

```
Flags:
      --preview     Preview the generated README locally (no push)
      --push        Push the generated README to your profile repo
      --top int     Number of top projects to display (default 5)

Global Flags:
      --token string      GitHub personal access token (or set GITHUB_TOKEN)
      --username string   GitHub username
```

### Interactive setup

```bash
githubautoprofil init
```

Walks you through setting up your username and token.

## GitHub Action

Add this workflow to your **profile repo** (`username/username`) to keep your README updated automatically:

```yaml
# .github/workflows/profile-readme.yml
name: Update Profile README

on:
  schedule:
    - cron: '0 0 * * 1' # Every Monday at midnight UTC
  workflow_dispatch:     # Manual trigger

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: yoanbernabeu/GitHubAutoProfil@v1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          username: ${{ github.repository_owner }}
          top: '5'
```

### Action inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `github_token` | Yes | — | GitHub token with `repo` scope |
| `username` | Yes | — | GitHub username |
| `top` | No | `5` | Number of top projects to display |
| `version` | No | `latest` | GitHubAutoProfil version to use |

## How it works

```
GitHub API ──→ Fetch Profile ──→ Fetch Repos ──→ Fetch Languages
                                                        │
                  README.md ◀── Render Template ◀── Analyze
                      │              │
                   Push to       Top repos by ⭐
                username/       Recent by push date
                username       Languages (normalized)
```

1. **Fetch** — pulls your profile, all owned repos (excluding forks/archived), and language data via the GitHub API
2. **Analyze** — ranks repos by stars and recent activity, aggregates languages with per-repo normalization (each repo weighs equally, regardless of size)
3. **Render** — generates Markdown from an embedded Go template with Shields.io badges
4. **Push** — creates or updates `README.md` in your `username/username` repo

## Example output

See it live: [github.com/yoanbernabeu](https://github.com/yoanbernabeu)

## Token

- **Preview mode** (`--preview`): no token needed — uses the public GitHub API (60 requests/hour)
- **Push mode** (`--push`): requires a [Personal Access Token](https://github.com/settings/tokens) with `repo` scope (5,000 requests/hour)

You can pass the token via `--token` flag or the `GITHUB_TOKEN` environment variable.

## License

MIT

---

<p align="center">
  <sub>Built with ❤️ by <a href="https://github.com/yoanbernabeu">Yoan Bernabeu</a></sub>
</p>
