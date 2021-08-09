package controller

import (
	"errors"
	"os"

	"github.com/m-mizutani/code-scanning-notify/pkg/domain/interfaces"
	"github.com/m-mizutani/code-scanning-notify/pkg/domain/model"
	"github.com/m-mizutani/code-scanning-notify/pkg/usecase"
	"github.com/m-mizutani/code-scanning-notify/pkg/utils"
	"github.com/m-mizutani/goerr"
	cli "github.com/urfave/cli/v2"
)

type Controller struct {
	usecase interfaces.Usecase
}

func New() *Controller {
	return &Controller{
		usecase: usecase.New(),
	}
}

type cliConfig struct {
	model.InputNotify
	model.Config
	logLevel string
}

func (x *Controller) CLI(args []string) {
	var cliCfg cliConfig

	app := &cli.App{
		Name:  "code-scanning-notify",
		Usage: "GitHub Action to notify Code Scanning Security Alert to PR",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "source-branch",
				Usage:       "Source branch",
				Aliases:     []string{"s"},
				Destination: &cliCfg.Source,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "target-branch",
				Usage:       "Target branch",
				Aliases:     []string{"t"},
				Destination: &cliCfg.Target,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "Repository owner/repo_name format",
				Aliases:     []string{"r"},
				Destination: &cliCfg.Repo,
				Required:    true,
			},
			&cli.Int64Flag{
				Name:        "issue-id",
				Usage:       "Issue (PR) ID",
				Aliases:     []string{"i"},
				Destination: &cliCfg.IssueID,
				Required:    true,
			},

			&cli.StringFlag{
				Name:        "github-token",
				Usage:       "GitHub Token",
				EnvVars:     []string{"GITHUB_TOKEN"},
				Destination: &cliCfg.GitHubToken,
			},

			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "Log level",
				Aliases:     []string{"l"},
				Destination: &cliCfg.logLevel,
				Value:       "info",
			},
		},

		Before: func(c *cli.Context) error {
			if err := utils.SetLogLevel(cliCfg.logLevel); err != nil {
				return err
			}
			return nil
		},

		Action: func(c *cli.Context) error {
			utils.Logger.Info().Interface("notify", cliCfg.InputNotify).Msg("Starting")

			x.usecase.SetConfig(&cliCfg.Config)
			if err := x.usecase.Notify(cliCfg.InputNotify); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		ev := utils.Logger.Error()

		var goErr *goerr.Error
		if errors.As(err, &goErr) {
			for k, v := range goErr.Values() {
				ev = ev.Interface(k, v)
			}
		}
		ev.Msg(err.Error())
	}
}
