package model

import "github.com/m-mizutani/goerr"

var (
	ErrGitHubAPI       = goerr.New("github API error")
	ErrInvalidRepoName = goerr.New("invalid repository name")
)
