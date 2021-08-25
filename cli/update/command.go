package update

import (
	"github.com/ccremer/greposync/cli/flags"
	"github.com/urfave/cli/v2"
)

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.cliCommand
}

func (c *Command) createCliCommand() *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Action: c.runCommand,
		Before: c.validateUpdateCommand,
		Flags: flags.CombineWithGlobalFlags(
			&cli.StringFlag{
				Name:    dryRunFlagName,
				Aliases: []string{"d"},
				Usage:   "Select a dry run mode. Allowed values: offline (do not run any Git commands), commit (commit, but don't push), push (push, but don't touch PRs)",
			},
			&cli.BoolFlag{
				Name:  amendFlagName,
				Usage: "Amend previous commit.",
			},
			&cli.BoolFlag{
				Name:  prCreateFlagName,
				Usage: "Create a PullRequest on a supported git hoster after pushing to remote.",
			},
			&cli.StringFlag{
				Name:  prBodyFlagName,
				Usage: "Markdown-enabled body of the PullRequest. It will load from an existing file if this is a path. Content can be templated. Defaults to commit message.",
			},
		),
	}
}
