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
	// ConfigDefaultName is the fallback file name of the YAML file containing the default template values.
	ConfigDefaultName = "config_defaults.yml"
	// GrepoSyncFileName is the default file name of the YAML file containing the main settings.
	GrepoSyncFileName = "greposync.yml"

	logLevelFlagName    = "log-level"
	projectRootFlagName = "project-root"
	projectJobsFlagName = "project-jobs"
)

func CreateCLI(version, commit, date string) {
	dateLayout := "2006-01-02"
	t, err := time.Parse(dateLayout, date)
	printer.CheckIfError(err)

	config = cfg.NewDefaultConfig()
	globalFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   fmt.Sprintf("Shorthand for --%s=debug", logLevelFlagName),
		},
		&cli.StringFlag{
			Name:  logLevelFlagName,
			Usage: "Log level. Allowed values are [debug, info, warn, error].",
			Value: config.Log.Level,
		},
		&cli.PathFlag{
			Name:  projectRootFlagName,
			Usage: "Local directory path where git clones repositories into.",
			Value: config.Project.RootDir,
		},
		&cli.IntFlag{
			Name:    projectJobsFlagName,
			Usage:   "Jobs is the number of parallel jobs to run. 1 basically means that jobs are run in sequence.",
			Aliases: []string{"j"},
			Value:   1,
		},
	}
	updateCommand := NewUpdateCommand(config)
	initCommand := NewInitCommand()
	app = &cli.App{
		Name:                 "greposync",
		Usage:                "git-repo-sync: Shameless reimplementation of ModuleSync in Go",
		Version:              fmt.Sprintf("%s, commit %s, date %s", version, commit[0:7], t.Format(dateLayout)),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			initCommand.createInitCommand(),
			updateCommand.createUpdateCommand(),
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