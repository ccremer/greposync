package application

import (
	"errors"
	"os"

	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/application/initialize"
	"github.com/ccremer/greposync/application/labels"
	"github.com/ccremer/greposync/application/update"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type (
	App struct {
		app    *cli.App
		config *cfg.Configuration
		log    logr.Logger
	}
)

// NewApp initializes the CLI application.
func NewApp(info VersionInfo, config *cfg.Configuration,
	labelCommand *labels.Command,
	updateCommand *update.Command,
	initializeCommand *initialize.Command,
	factory logging.LoggerFactory,
) *App {
	flags.InitGlobalFlags(config)
	app := &App{
		log:    factory.NewGenericLogger(""),
		config: config,
	}
	a := &cli.App{
		Name:                 "greposync",
		Usage:                "git-repo-sync: Shameless reimplementation of ModuleSync in Go",
		Version:              info.String(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			initializeCommand.GetCliCommand(),
			labelCommand.GetCliCommand(),
			updateCommand.GetCliCommand(),
		},
		ExitErrHandler: clierror.NewErrorHandler(app.log),
	}
	app.app = a
	return app
}

// Run the CLI application
func (a *App) Run() {
	err := a.app.Run(os.Args)
	if err != nil {
		if errors.Is(err, clierror.ErrPipeline) {
			// Ignore pipeline errors as they are printed separately
			a.log.Error(nil, "An error occurred in one of the repositories. See the log printed above")
		} else {
			a.log.Error(nil, err.Error())
		}
		os.Exit(1)
	}
}
