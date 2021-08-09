package infra

import "github.com/m-mizutani/cs-alert-notify/pkg/domain/interfaces"

func New() *interfaces.Factories {
	return &interfaces.Factories{
		NewGitHub: NewGitHubClient,
	}
}
