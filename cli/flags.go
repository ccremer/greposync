package cli

import "github.com/urfave/cli/v2"

const (
	projectIncludeFlagName = "project-include"
	projectExcludeFlagName = "project-exclude"
)

var (
	projectIncludeFlag = &cli.StringFlag{
		Name:  projectIncludeFlagName,
		Usage: "Includes only repositories in the update that match the given filter (regex).",
	}
	projectExcludeFlag = &cli.StringFlag{
		Name:  projectExcludeFlagName,
		Usage: "Excludes repositories from updating that match the given filter (regex). Repositories matching both include and exclude filter are still excluded.",
	}
)
