package labels

import (
	"github.com/ccremer/greposync/application/flags"
	"github.com/urfave/cli/v2"
)

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.createCommand()
}

func (c *Command) createCommand() *cli.Command {
	cFlags := []cli.Flag{
		flags.NewLogLevelFlag(&c.cfg.Log.Level),
		flags.NewShowLogFlag(&c.cfg.Log.ShowLog),

		flags.NewJobsFlag(&c.cfg.Project.Jobs),
		flags.NewSkipBrokenFlag(&c.cfg.Project.SkipBroken),
		flags.NewIncludeFlag(&c.cfg.Project.Include),
		flags.NewExcludeFlag(&c.cfg.Project.Exclude),

		flags.NewGitCommitBranchFlag(&c.cfg.Git.CommitBranch),
		flags.NewGitDefaultNamespaceFlag(&c.appService.repoStore.DefaultNamespace),
		flags.NewGitRootDirFlag(&c.appService.repoStore.ParentDir),
	}
	return &cli.Command{
		Name:   "labels",
		Usage:  "Synchronizes repository labels",
		Before: flags.And(flags.FromYAML(cFlags), c.validateCommand),
		Action: c.runCommand,
		Flags:  cFlags,
	}
}
