package cli

import "github.com/urfave/cli/v2"

const (
	ProjectIncludeFlagName = "project-include"
	ProjectExcludeFlagName = "project-exclude"
	projectRootFlagName    = "project-root"
	ProjectJobsFlagName    = "project-jobs"
)

var (
	JobsMinimumCount = 1
	JobsMaximumCount = 8
)

func NewProjectIncludeFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  ProjectIncludeFlagName,
		Usage: "Includes only repositories in the update that match the given filter (regex).",
	}
}

func NewProjectExcludeFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  ProjectExcludeFlagName,
		Usage: "Excludes repositories from updating that match the given filter (regex). Repositories matching both include and exclude filter are still excluded.",
	}
}
