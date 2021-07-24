package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/initialize"
	"github.com/ccremer/greposync/cli/labels"
	"github.com/ccremer/greposync/cli/update"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

var (
	app         *cli.App
	config      *cfg.Configuration
	// ConfigDefaultName is the fallback file name of the YAML file containing the default template values.
	ConfigDefaultName = "config_defaults.yml"
	// GrepoSyncFileName is the default file name of the YAML file containing the main settings.
	GrepoSyncFileName = "greposync.yml"
)

func CreateCLI(version, commit, date string) {
	dateLayout := "2006-01-02"
	t, err := time.Parse(dateLayout, date)
	printer.CheckIfError(err)

	config = cfg.NewDefaultConfig()
	app = &cli.App{
		Name:                 "greposync",
		Usage:                "git-repo-sync: Shameless reimplementation of ModuleSync in Go",
		Version:              fmt.Sprintf("%s, commit %s, date %s", version, commit[0:7], t.Format(dateLayout)),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			initialize.NewCommand(config).GetCliCommand(),
			labels.NewCommand(config).GetCliCommand(),
			update.NewCommand(config).GetCliCommand(),
		},
		Compiled: t,
		ExitErrHandler: func(context *cli.Context, err error) {
			_ = cli.ShowCommandHelp(context, context.Command.Name)
		},
	}
}

// Run the CLI application
func Run() {
	err := app.Run(os.Args)
	printer.CheckIfError(err)
}
