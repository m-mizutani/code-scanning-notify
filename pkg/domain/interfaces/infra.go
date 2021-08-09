package interfaces

import "github.com/google/go-github/v37/github"

type GitHub interface {
	ListAlertsForRepo(owner, repo string, opts *github.AlertListOptions) ([]*github.Alert, error)
	CreateComment(owner, repo string, number int, comment *github.IssueComment) error
}

type GitHubFactory func(token string) GitHub

type Factories struct {
	NewGitHub GitHubFactory
}
