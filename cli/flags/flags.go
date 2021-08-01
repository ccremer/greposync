package flags

import (
	"fmt"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

const (
	// ProjectIncludeFlagName is the name on the CLI
	ProjectIncludeFlagName = "project-include"
	// ProjectExcludeFlagName is the name on the CLI
	ProjectExcludeFlagName = "project-exclude"
	// ProjectJobsFlagName is the name on the CLI
	ProjectJobsFlagName = "project-jobs"

	projectRootFlagName = "project-root"
	logLevelFlagName    = "log-level"
)

var (
	JobsMinimumCount = 1
	JobsMaximumCount = 8
	globalFlags      []cli.Flag
)

func InitGlobalFlags(config *cfg.Configuration) []cli.Flag {

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
			Name:    ProjectJobsFlagName,
			Usage:   "Jobs is the number of parallel jobs to run. 1 basically means that jobs are run in sequence.",
			Aliases: []string{"j"},
			Value:   1,
		},
		NewProjectIncludeFlag(),
		NewProjectExcludeFlag(),
	}
	return globalFlags
}

// NewProjectIncludeFlag returns a new flag that is meant to include Git repositories.
func NewProjectIncludeFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  ProjectIncludeFlagName,
		Usage: "Includes only repositories in the update that match the given filter (regex).",
	}
}

// NewProjectExcludeFlag returns a new flag that is meant to exclude Git repositories.
func NewProjectExcludeFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  ProjectExcludeFlagName,
		Usage: "Excludes repositories from updating that match the given filter (regex). Repositories matching both include and exclude filter are still excluded.",
	}
}

// CombineWithGlobalFlags combines the given flags with the global flags.
// The given flags are appended, so the global flags are first in the list.
func CombineWithGlobalFlags(flags ...cli.Flag) []cli.Flag {
	for _, flag := range flags {
		globalFlags = append(globalFlags, flag)
	}
	return globalFlags
}

// ValidateGlobalFlags validates whether the global flags are given as expected.
func ValidateGlobalFlags(ctx *cli.Context, config *cfg.Configuration) error {
	if ctx.Bool("verbose") {
		config.Log.Level = "debug"
		printer.DefaultLevel = printer.LevelDebug
	}
	return nil
}
