package update

import (
	"context"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
	"github.com/go-logr/logr"
)

type updatePipeline struct {
	log        logr.Logger
	repo       *domain.GitRepository
	appService *AppService
}

func (c *updatePipeline) clone(_ context.Context) error {
	return c.appService.repoStore.Clone(c.repo)
}

func (c *updatePipeline) fetch(_ context.Context) error {
	return c.appService.repoStore.Fetch(c.repo)
}

func (c *updatePipeline) pull(_ context.Context) error {
	return c.appService.repoStore.Pull(c.repo)
}

func (c *updatePipeline) checkout(_ context.Context) error {
	return c.appService.repoStore.Checkout(c.repo)
}

func (c *updatePipeline) reset(_ context.Context) error {
	return c.appService.repoStore.Reset(c.repo)
}

func (c *updatePipeline) add(_ context.Context) error {
	return c.appService.repoStore.Add(c.repo)
}

func (c *updatePipeline) commit(_ context.Context) error {
	err := c.appService.repoStore.Commit(c.repo, domain.CommitOptions{
		Message: c.appService.cfg.Git.CommitMessage,
		Amend:   c.appService.cfg.Git.Amend,
	})
	return err
}

func (c *updatePipeline) diff(_ context.Context) error {
	diff, err := c.appService.repoStore.Diff(c.repo, domain.DiffOptions{
		WorkDirToHEAD: c.appService.cfg.Git.SkipCommit, // If we don't commit, show the unstaged changes
	})
	if err != nil {
		return err
	}
	c.appService.diffPrinter.PrintDiff("Diff: "+c.repo.URL.GetFullName(), diff)
	return nil
}

func (c *updatePipeline) push(_ context.Context) error {
	err := c.appService.repoStore.Push(c.repo, domain.PushOptions{
		Force: c.appService.cfg.Git.ForcePush,
	})
	return err
}

func (c *updatePipeline) renderTemplates(_ context.Context) error {
	err := c.appService.renderService.RenderTemplates(domain.RenderContext{
		Repository:    c.repo,
		ValueStore:    c.appService.valueStore,
		TemplateStore: c.appService.templateStore,
		Engine:        c.appService.engine,
	})
	return err
}

func (c *updatePipeline) cleanupUnwantedFiles(_ context.Context) error {
	err := c.appService.cleanupService.CleanupUnwantedFiles(domain.CleanupContext{
		Repository: c.repo,
		ValueStore: c.appService.valueStore,
	})
	return err
}

func (c *updatePipeline) ensurePullRequest(_ context.Context) error {
	if c.repo.PullRequest == nil {
		err := c.appService.prService.NewPullRequestForRepository(domain.PullRequestServiceContext{
			Repository:     c.repo,
			TemplateEngine: c.appService.engine,
			Body:           c.appService.cfg.PullRequest.BodyTemplate,
			Title:          c.appService.cfg.PullRequest.Subject,
			TargetBranch:   c.repo.DefaultBranch,
		})
		if err != nil {
			return err
		}
	}
	if err := c.repo.PullRequest.AttachLabels(domain.FromStringSlice(c.appService.cfg.PullRequest.Labels)); err != nil {
		return err
	}
	err := c.appService.prStore.EnsurePullRequest(c.repo)
	return err
}

func (c *updatePipeline) dirMissing() pipeline.Predicate {
	return func(_ context.Context) bool {
		return !c.repo.RootDir.DirExists()
	}
}

func (c *updatePipeline) isDirty() pipeline.Predicate {
	return func(_ context.Context) bool {
		return c.appService.repoStore.IsDirty(c.repo)
	}
}

func (c *updatePipeline) hasCommits() pipeline.Predicate {
	return func(_ context.Context) bool {
		return true
		// TODO: There should be a better way to determine whether to push...

		/*		hasCommits, err := repositorystore.HasCommitsBetween(c.repo, c.repo.DefaultBranch, c.repo.CommitBranch)
				if err != nil {
					c.log.WarnF("%w", err)
				}
				return hasCommits*/
	}
}

func (c *updatePipeline) fetchPullRequest(_ context.Context) error {
	pr, err := c.appService.prStore.FindMatchingPullRequest(c.repo)
	c.repo.PullRequest = pr
	return err
}
