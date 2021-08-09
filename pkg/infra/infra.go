package infra

import "github.com/m-mizutani/code-scanning-notify/pkg/domain/interfaces"

func New() *interfaces.Factories {
	return &interfaces.Factories{
		NewGitHub: NewGitHubClient,
	}
}
