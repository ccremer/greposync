package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/clierror"
	"github.com/ccremer/greposync/cli/flags"
	"github.com/ccremer/greposync/cli/initialize"
	"github.com/ccremer/greposync/cli/labels"
	"github.com/ccremer/greposync/cli/update"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

type (
	App struct {
		app    *cli.App
		config *cfg.Configuration
	}
	VersionInfo struct {
		Version string
		Commit  string
		Date    string
	}
)

// NewApp initializes the CLI application.
func NewApp(info VersionInfo, config *cfg.Configuration) *App {
	dateLayout := "2006-01-02"
	t, err := time.Parse(dateLayout, info.Date)
	if err != nil {
		printer.DefaultPrinter.ErrorF(err.Error())
		os.Exit(2)
	}

	flags.InitGlobalFlags(config)
	a := &cli.App{
		Name:                 "greposync",
		Usage:                "git-repo-sync: Shameless reimplementation of ModuleSync in Go",
		Version:              fmt.Sprintf("%s, commit %s, date %s", info.Version, info.Commit[0:7], t.Format(dateLayout)),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			initialize.NewCommand(config).GetCliCommand(),
			labels.NewCommand(config).GetCliCommand(),
			update.NewCommand(config).GetCliCommand(),
		},
		Compiled:       t,
		ExitErrHandler: clierror.ErrorHandler,
	}
	return &App{
		app:    a,
		config: config,
	}
}

// Run the CLI application
func (a *App) Run() {
	err := a.app.Run(os.Args)
	if err != nil {
		printer.DefaultPrinter.ErrorF(err.Error())
		os.Exit(1)
	}
}
