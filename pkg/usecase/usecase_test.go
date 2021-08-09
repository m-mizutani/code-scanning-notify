package usecase_test

import (
	"testing"

	"github.com/google/go-github/v37/github"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/interfaces"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/model"
	"github.com/m-mizutani/cs-alert-notify/pkg/infra"
	"github.com/m-mizutani/cs-alert-notify/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotify(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		uc := usecase.New()
		githubMock, githubNew := infra.NewGitHubMock()
		uc.InjectFactories(&interfaces.Factories{
			NewGitHub: githubNew,
		})

		countListAlertsForRepoMock := 0
		githubMock.ListAlertsForRepoMock = func(owner, repo string, opts *github.AlertListOptions) ([]*github.Alert, error) {
			countListAlertsForRepoMock++

			assert.Equal(t, "clock", owner)
			assert.Equal(t, "tower", repo)
			switch opts.Ref {
			case "heads/dev":
				return []*github.Alert{
					{
						Rule: &github.Rule{
							ID:          github.String("blue"),
							Description: github.String("blue"),
						},
						HTMLURL: github.String("https://example.com/1"),
						MostRecentInstance: &github.MostRecentInstance{
							Message: &github.Message{
								Text: github.String("sample1"),
							},
						},
					},
					{
						Rule: &github.Rule{
							ID:          github.String("red"),
							Description: github.String("red"),
						},
						HTMLURL: github.String("https://example.com/2"),
						MostRecentInstance: &github.MostRecentInstance{
							Message: &github.Message{
								Text: github.String("sample2"),
							},
						},
					},
				}, nil

			case "heads/release":
				return []*github.Alert{
					{
						Rule: &github.Rule{
							ID:          github.String("orange"),
							Description: github.String("orange"),
						},
						HTMLURL: github.String("https://example.com/3"),
						MostRecentInstance: &github.MostRecentInstance{
							Message: &github.Message{
								Text: github.String("sample3"),
							},
						},
					},
					{
						Rule: &github.Rule{
							ID:          github.String("red"),
							Description: github.String("red"),
						},
						HTMLURL: github.String("https://example.com/2"),
						MostRecentInstance: &github.MostRecentInstance{
							Message: &github.Message{
								Text: github.String("sample2"),
							},
						},
					},
				}, nil
			}
			require.FailNow(t, "invalid Ref")
			return nil, nil
		}

		countCreatedComment := 0
		githubMock.CreateCommentMock = func(owner, repo string, number int, comment *github.IssueComment) error {
			assert.Equal(t, "clock", owner)
			assert.Equal(t, "tower", repo)
			assert.Equal(t, 5, number)

			require.NotNil(t, comment.Body)
			assert.Contains(t, *comment.Body, "🚨")
			assert.Contains(t, *comment.Body, "✅")
			assert.Contains(t, *comment.Body, "⚠")

			countCreatedComment++
			return nil
		}

		require.NoError(t, uc.Notify(model.InputNotify{
			Source:  "heads/dev",
			Target:  "heads/release",
			Repo:    "clock/tower",
			IssueID: 5,
		}))
		assert.Equal(t, 1, countCreatedComment)
		assert.Equal(t, 2, countListAlertsForRepoMock)
	})
}
