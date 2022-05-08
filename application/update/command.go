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
		flags.NewIncludeFlag(&c.appService.repoStore.IncludeFilter),
		flags.NewExcludeFlag(&c.appService.repoStore.ExcludeFilter),
		flags.NewDryRunFlag(&c.dryRunFlag),

		flags.NewGitRootDirFlag(&c.appService.repoStore.ParentDir),
		flags.NewGitAmendFlag(&c.cfg.Git.Amend),
		flags.NewGitForcePushFlag(&c.cfg.Git.ForcePush),
		flags.NewGitCommitBranchFlag(&c.appService.repoStore.CommitBranch),
		flags.NewGitDefaultNamespaceFlag(&c.appService.repoStore.DefaultNamespace),
		flags.NewGitCommitMessageFlag(&c.cfg.Git.CommitMessage),
		flags.NewGitBaseURLFlag(&c.appService.repoStore.BaseURL),

		flags.NewPRCreateFlag(&c.cfg.PullRequest.Create),
		flags.NewPRBodyFlag(&c.cfg.PullRequest.BodyTemplate),
		flags.NewPRSubjectFlag(&c.cfg.PullRequest.Subject),
		flags.NewPRTargetBranchFlag(&c.cfg.PullRequest.TargetBranch),
		flags.NewPRLabelsFlag(&c.PrLabels),

		flags.NewTemplateRootDirFlag(&c.appService.templateStore.RootDir),
	}
	return &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Before: flags.And(flags.FromYAML(cFlags), c.validateUpdateCommand),
		Action: c.runCommand,
		Flags:  cFlags,
	}
}
