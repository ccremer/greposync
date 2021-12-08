package update

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/domain"
	"github.com/go-logr/logr"
)

type pipelineContext struct {
	log        logr.Logger
	repo       *domain.GitRepository
	appService *AppService
}

func (c *pipelineContext) clone() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Clone)
}

func (c *pipelineContext) fetch() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Fetch)
}

func (c *pipelineContext) pull() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Pull)
}

func (c *pipelineContext) checkout() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Checkout)
}

func (c *pipelineContext) reset() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Reset)
}

func (c *pipelineContext) add() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Add)
}

func (c *pipelineContext) commit() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		err := c.appService.repoStore.Commit(c.repo, domain.CommitOptions{
			Message: c.appService.cfg.Git.CommitMessage,
			Amend:   c.appService.cfg.Git.Amend,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) diff() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		diff, err := c.appService.repoStore.Diff(c.repo, domain.DiffOptions{
			WorkDirToHEAD: c.appService.cfg.Git.SkipCommit, // If we don't commit, show the unstaged changes
		})
		if err != nil {
			return pipeline.Result{Err: err}
		}
		c.appService.diffPrinter.PrintDiff("Diff: "+c.repo.URL.GetFullName(), diff)
		return pipeline.Result{}
	}
}

func (c *pipelineContext) push() pipeline.ActionFunc {
	return func(ctx pipeline.Context) pipeline.Result {
		err := c.appService.repoStore.Push(c.repo, domain.PushOptions{
			Force: c.appService.cfg.Git.ForcePush,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) renderTemplates() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		err := c.appService.renderService.RenderTemplates(domain.RenderContext{
			Repository:    c.repo,
			ValueStore:    c.appService.valueStore,
			TemplateStore: c.appService.templateStore,
			Engine:        c.appService.engine,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) cleanupUnwantedFiles() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		err := c.appService.cleanupService.CleanupUnwantedFiles(domain.CleanupContext{
			Repository: c.repo,
			ValueStore: c.appService.valueStore,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) ensurePullRequest() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		if c.repo.PullRequest == nil {
			err := c.appService.prService.NewPullRequestForRepository(domain.PullRequestServiceContext{
				Repository:     c.repo,
				TemplateEngine: c.appService.engine,
				Body:           c.appService.cfg.PullRequest.BodyTemplate,
				Title:          c.appService.cfg.PullRequest.Subject,
				TargetBranch:   c.repo.DefaultBranch,
			})
			if err != nil {
				return pipeline.Result{Err: err}
			}
		}
		if err := c.repo.PullRequest.AttachLabels(domain.FromStringSlice(c.appService.cfg.PullRequest.Labels)); err != nil {
			return pipeline.Result{Err: err}
		}
		err := c.appService.prStore.EnsurePullRequest(c.repo)
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) dirMissing() predicate.Predicate {
	return func(ctx pipeline.Context, step pipeline.Step) bool {
		return !c.repo.RootDir.DirExists()
	}
}

func (c *pipelineContext) isDirty() predicate.Predicate {
	return func(ctx pipeline.Context, step pipeline.Step) bool {
		return c.appService.repoStore.IsDirty(c.repo)
	}
}

func (c *pipelineContext) hasCommits() predicate.Predicate {
	return func(ctx pipeline.Context, step pipeline.Step) bool {
		return true
		// TODO: There should be a better way to determine whether to push...

		/*		hasCommits, err := repositorystore.HasCommitsBetween(c.repo, c.repo.DefaultBranch, c.repo.CommitBranch)
				if err != nil {
					c.log.WarnF("%w", err)
				}
				return hasCommits*/
	}
}

func (c *pipelineContext) toAction(f func(repository *domain.GitRepository) error) pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		err := f(c.repo)
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) fetchPullRequest() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		pr, err := c.appService.prStore.FindMatchingPullRequest(c.repo)
		c.repo.PullRequest = pr
		return pipeline.Result{Err: err}
	}
}
