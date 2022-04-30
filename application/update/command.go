package update

import (
	"github.com/ccremer/greposync/application/flags"
	"github.com/urfave/cli/v2"
)

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.createCliCommand()
}

func (c *Command) createCliCommand() *cli.Command {
	cFlags := []cli.Flag{
		flags.NewLogLevelFlag(&c.cfg.Log.Level),
		flags.NewShowLogFlag(&c.cfg.Log.ShowLog),
		flags.NewShowDiffFlag(&c.cfg.Log.ShowDiff),

		flags.NewJobsFlag(&c.cfg.Project.Jobs),
		flags.NewSkipBrokenFlag(&c.cfg.Project.SkipBroken),
		flags.NewIncludeFlag(&c.cfg.Project.Include),
		flags.NewExcludeFlag(&c.cfg.Project.Exclude),
		flags.NewDryRunFlag(&c.dryRunFlag),

		flags.NewGitRootDirFlag(&c.cfg.Project.RootDir),
		flags.NewGitAmendFlag(&c.cfg.Git.Amend),
		flags.NewGitForcePushFlag(&c.cfg.Git.ForcePush),
		flags.NewGitCommitBranchFlag(&c.cfg.Git.CommitBranch),
		flags.NewGitDefaultNamespaceFlag(&c.cfg.Git.Namespace),
		flags.NewGitCommitMessageFlag(&c.cfg.Git.CommitMessage),

		flags.NewPRCreateFlag(&c.cfg.PullRequest.Create),
		flags.NewPRBodyFlag(&c.cfg.PullRequest.BodyTemplate),
		flags.NewPRSubjectFlag(&c.cfg.PullRequest.Subject),
		flags.NewPRTargetBranchFlag(&c.cfg.PullRequest.TargetBranch),
		flags.NewPRLabelsFlag(&c.PrLabels),

		flags.NewTemplateRootDirFlag(&c.cfg.Template.RootDir),
	}
	return &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Before: flags.And(flags.FromYAML(cFlags), c.validateUpdateCommand),
		Action: c.runCommand,
		Flags:  cFlags,
	}
}
