package usecase

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v37/github"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/interfaces"
	"github.com/m-mizutani/cs-alert-notify/pkg/domain/model"
	"github.com/m-mizutani/cs-alert-notify/pkg/infra"
	"github.com/m-mizutani/cs-alert-notify/pkg/utils"
	"github.com/m-mizutani/goerr"
)

type Usecase struct {
	config    *model.Config
	factories *interfaces.Factories
}

func New() *Usecase {
	return &Usecase{
		config:    &model.Config{},
		factories: infra.New(),
	}
}

var logger = utils.Logger

func (x *Usecase) Notify(input model.InputNotify) error {
	githubClient := x.factories.NewGitHub(x.config.GitHubToken)

	repo := strings.Split(input.Repo, "/")
	if len(repo) != 2 {
		return goerr.Wrap(model.ErrInvalidRepoName).With("repo", input.Repo)
	}

	srcAlerts, err := githubClient.ListAlertsForRepo(repo[0], repo[1], &github.AlertListOptions{
		Ref: input.Source,
	})
	if err != nil {
		return err
	}

	tgtAlerts, err := githubClient.ListAlertsForRepo(repo[0], repo[1], &github.AlertListOptions{
		Ref: input.Target,
	})
	if err != nil {
		return err
	}

	diff := diffAlert(srcAlerts, tgtAlerts)
	if diff.NoAlert() {
		logger.Info().Msg("No added/fixed alert, quiting")
		return nil
	}

	comment := &github.IssueComment{
		Body: github.String(buildComment(repo[0], repo[1], diff)),
	}
	if err := githubClient.CreateComment(repo[0], repo[1], int(input.IssueID), comment); err != nil {
		return err
	}

	return nil
}

type alertDiff struct {
	Added    []*github.Alert
	Remained []*github.Alert
	Deleted  []*github.Alert
}

func (x *alertDiff) NoAlert() bool {
	return len(x.Added) == 0 && len(x.Remained) == 0 && len(x.Deleted) == 0
}

func diffAlert(src, tgt []*github.Alert) *alertDiff {
	var diff alertDiff
	makeAlertMap := func(alerts []*github.Alert) map[int64]*github.Alert {
		resp := make(map[int64]*github.Alert)
		for _, alert := range alerts {
			resp[alert.ID()] = alert
		}
		return resp
	}

	srcMap := makeAlertMap(src)
	tgtMap := makeAlertMap(tgt)

	for id, alert := range srcMap {
		if _, ok := tgtMap[id]; !ok {
			diff.Added = append(diff.Added, alert)
		} else {
			diff.Remained = append(diff.Remained, alert)
		}
	}
	for id, alert := range tgtMap {
		if _, ok := srcMap[id]; !ok {
			diff.Deleted = append(diff.Deleted, alert)
		}
	}

	return &diff
}

func buildComment(owner, repo string, diff *alertDiff) string {
	var msg string

	str := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	if len(diff.Added) > 0 {
		msg += "### 🚨 New alerts\n"
		for _, alert := range diff.Added {
			msg += fmt.Sprintf("- [%s](%s): %s\n",
				str(alert.Rule.Description),
				str(alert.HTMLURL),
				alert.MostRecentInstance.Message.GetText())
		}
		msg += "\n"
	}

	if len(diff.Deleted) > 0 {
		msg += "### ✅ Fixed alerts\n"
		for _, alert := range diff.Deleted {
			msg += "- " + *alert.Rule.Description + "\n"
		}
		msg += "\n"
	}

	if len(diff.Remained) > 0 {
		msg += "### ⚠ Remained alerts\n"
		msg += fmt.Sprintf("Remained %d alerts. See list from [here](https://github.com/%s/%s/security/code-scanning).\n", len(diff.Remained), owner, repo)
	}

	return msg
}

func (x *Usecase) SetConfig(cfg *model.Config) {
	x.config = cfg
}
