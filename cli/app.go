package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

var (
	app         *cli.App
	globalFlags []cli.Flag
	config      *cfg.Configuration
)

func CreateCLI(version, commit, date string) {
	dateLayout := "2006-01-02"
	t, err := time.Parse(dateLayout, date)
	printer.CheckIfError(err)

	c := cfg.NewDefaultConfig()
	globalFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Shorthand for --log-level=debug",
		},
		&cli.StringFlag{
			Name:        "log-level",
			Destination: &c.Log.Level,
			Usage:       "Log level. Allowed values are [debug, info, warn, error].",
			Value:       "info",
		},
	}
	app = &cli.App{
		Name:                 "Git-Repo-Sync",
		Usage:                "Shameless reimplementation of ModuleSync in Go",
		Version:              fmt.Sprintf("%s, commit %s, date %s", version, commit[0:7], t.Format(dateLayout)),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			createUpdateCommand(c),
		},
		Compiled: t,
		ExitErrHandler: func(context *cli.Context, err error) {
			_ = cli.ShowCommandHelp(context, context.Command.Name)
		},
		Before: func(context *cli.Context) error {
			return cfg.ParseConfig("gitreposync.yml", c)
		},
	}
	config = c
}

// Run the CLI application
func Run() {
	err := app.Run(os.Args)
	printer.CheckIfError(err)
}

func combineWithGlobalFlags(flags ...cli.Flag) []cli.Flag {
	for _, flag := range flags {
		globalFlags = append(globalFlags, flag)
	}
	return globalFlags
}

func validateGlobalFlags(ctx *cli.Context) error {
	if ctx.Bool("verbose") {
		config.Log.Level = "debug"
		printer.DefaultLevel = printer.LevelDebug
	}
	return nil
}
