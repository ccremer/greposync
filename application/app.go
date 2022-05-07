package application

import (
	"errors"
	"os"
	"sort"

	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/application/initialize"
	"github.com/ccremer/greposync/application/labels"
	"github.com/ccremer/greposync/application/test"
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

func init() {
	// Remove -v short option
	cli.VersionFlag.(*cli.BoolFlag).Aliases = nil
}

// NewApp initializes the CLI application.
func NewApp(info VersionInfo, config *cfg.Configuration,
	labelCommand *labels.Command,
	updateCommand *update.Command,
	initializeCommand *initialize.Command,
	testCommand *test.Command,
	factory logging.LoggerFactory,
) *App {
	app := &App{
		log:    factory.NewGenericLogger(""),
		config: config,
	}
	a := &cli.App{
		Name:  "greposync",
		Usage: "Managed Git repositories in bulk",
		Description: `At the heart of greposync is a template.
The template exists of files that are being rendered with various input variables and ultimately committed to a Git repository.
greposync enables you to keep multiple Git repositories aligned with all the skeleton files that you need.

While services like GitHub offer the functionality of template repository, once you generated a new repository from the template it's not being updated anymore.
Over time you'll do changes to your CI/CD workflows or Makefiles and you want the changes in all your popular repositories. 
greposync does just that.`,
		Version:              info.String(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			initializeCommand.GetCliCommand(),
			labelCommand.GetCliCommand(),
			updateCommand.GetCliCommand(),
			testCommand.GetCliCommand(),
		},
		ExitErrHandler: clierror.NewErrorHandler(app.log),
	}
	for _, command := range a.Commands {
		sort.Sort(cli.FlagsByName(command.Flags))
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
