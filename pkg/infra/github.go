package infra

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v37/github"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/interfaces"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/model"
	"github.com/m-mizutani/cs-alert-notify/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
}

func NewGitHubClient(token string) interfaces.GitHub {
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return &GitHubClient{
		client: github.NewClient(tc),
	}
}

func (x *GitHubClient) ListAlertsForRepo(owner string, repo string, opts *github.AlertListOptions) ([]*github.Alert, error) {
	ctx := context.Background()
	alerts, resp, err := x.client.CodeScanning.ListAlertsForRepo(ctx, owner, repo, opts)
	if err != nil {
		return nil, goerr.Wrap(err, "ListAlertsForRepo").
			With("owner", owner).With("repo", repo).With("opt", opts)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, goerr.Wrap(model.ErrGitHubAPI, "ListAlertsForRepo").
			With("body", body).With("code", resp.StatusCode).With("status", resp.Status)
	}

	return alerts, nil
}

func (x *GitHubClient) CreateComment(owner string, repo string, number int, comment *github.IssueComment) error {
	ctx := context.Background()
	created, resp, err := x.client.Issues.CreateComment(ctx, owner, repo, number, comment)
	if err != nil {
		return goerr.Wrap(err, "CreateComment").With("owner", owner).With("repo", repo).With("comment", comment).With("id", number)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return goerr.Wrap(model.ErrGitHubAPI, "CreateComment").With("body", body).With("code", resp.StatusCode).With("status", resp.Status)
	}

	utils.Logger.Debug().Interface("comment", created).Msg("comment created")

	return nil
}

type GitHubMock struct {
	ListAlertsForRepoMock func(owner string, repo string, opts *github.AlertListOptions) ([]*github.Alert, error)
	CreateCommentMock     func(owner string, repo string, number int, comment *github.IssueComment) error
}

func NewGitHubMock() (*GitHubMock, interfaces.GitHubFactory) {
	mock := &GitHubMock{}
	return mock, func(token string) interfaces.GitHub { return mock }
}

func (x *GitHubMock) ListAlertsForRepo(owner string, repo string, opts *github.AlertListOptions) ([]*github.Alert, error) {
	return x.ListAlertsForRepoMock(owner, repo, opts)
}

func (x *GitHubMock) CreateComment(owner string, repo string, number int, comment *github.IssueComment) error {
	return x.CreateCommentMock(owner, repo, number, comment)
}
