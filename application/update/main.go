package update

import (
	"context"

	pipeline "github.com/ccremer/go-command-pipeline"
	instrumentation2 "github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/urfave/cli/v2"
)

const (
	dryRunFlagName    = "dry-run"
	prCreateFlagName  = "pr-create"
	prBodyFlagName    = "pr-bodyTemplate"
	amendFlagName     = "git-amend"
	forcePushFlagName = "git-forcePush"
	showDiffFlagName  = "log-showDiff"
)

type (
	// Command is a facade service for the update command that holds all dependent services and settings.
	Command struct {
		cfg             *cfg.Configuration
		repositories    []*domain.GitRepository
		appService      *AppService
		instrumentation instrumentation2.BatchInstrumentation
		logFactory      logging.LoggerFactory
	}
)

// NewCommand returns a new Command instance.
func NewCommand(
	cfg *cfg.Configuration,
	configurator *AppService,
	factory logging.LoggerFactory,
	instrumentation instrumentation2.BatchInstrumentation,
) *Command {
	c := &Command{
		cfg:             cfg,
		appService:      configurator,
		instrumentation: instrumentation,
		logFactory:      factory,
	}
	return c
}

func (c *Command) runCommand(_ *cli.Context) error {
	logger := c.logFactory.NewPipelineLogger("")
	p := pipeline.NewPipeline().AddBeforeHook(logger.Accept).WithSteps(
		pipeline.NewStepFromFunc("configure infrastructure", c.configureInfrastructure),
		pipeline.NewStepFromFunc("fetch managed repos config", c.fetchRepositories),
		pipeline.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.instrumentation.NewCollectErrorHandler(c.cfg.Project.SkipBroken)),
	)
	p.WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instrumentation.BatchPipelineCompleted(c.GetRepositories())
		return result.Err()
	})
	return p.Run().Err()
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {

	resetRepo := !c.cfg.Git.SkipReset
	enabledCommits := !c.cfg.Git.SkipCommit
	enabledPush := !c.cfg.Git.SkipPush
	showDiff := c.cfg.Log.ShowDiff
	createPR := c.cfg.PullRequest.Create

	repoCtx := &updatePipeline{
		repo:       r,
		appService: c.appService,
	}

	logger := c.logFactory.NewPipelineLogger(r.URL.GetFullName())
	p := pipeline.NewPipeline().AddBeforeHook(logger.Accept)
	p.WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ context.Context) error {
			c.instrumentation.PipelineForRepositoryStarted(repoCtx.repo)
			return nil
		}),

		pipeline.NewPipeline().AddBeforeHook(logger.Accept).
			WithNestedSteps("prepare workspace",

				pipeline.ToStep("clone repository", repoCtx.clone, repoCtx.dirMissing()),
				pipeline.ToStep("fetch", repoCtx.fetch, pipeline.Bool(resetRepo)),
				pipeline.ToStep("reset", repoCtx.reset, pipeline.Bool(resetRepo)),
				pipeline.ToStep("checkout branch", repoCtx.checkout, pipeline.Bool(resetRepo)),
				pipeline.ToStep("pull", repoCtx.pull, pipeline.Bool(resetRepo)),
			),

		pipeline.NewPipeline().AddBeforeHook(logger.Accept).
			WithNestedSteps("render",
				pipeline.NewStepFromFunc("render templates", repoCtx.renderTemplates),
				pipeline.NewStepFromFunc("cleanup unwanted files", repoCtx.cleanupUnwantedFiles),
			),

		pipeline.If(pipeline.And(pipeline.Bool(enabledCommits), repoCtx.isDirty()),
			pipeline.NewPipeline().
				AddBeforeHook(logger.Accept).
				WithNestedSteps("commit changes",
					pipeline.NewStepFromFunc("add", repoCtx.add),
					pipeline.NewStepFromFunc("commit", repoCtx.commit),
				),
		),

		pipeline.ToStep("show diff", repoCtx.diff, pipeline.Bool(showDiff)),
		pipeline.ToStep("push changes", repoCtx.push, pipeline.And(pipeline.Bool(enabledPush), repoCtx.hasCommits())),
		pipeline.ToStep("find existing pull request", repoCtx.fetchPullRequest, pipeline.Bool(createPR)),
		pipeline.ToStep("ensure pull request", repoCtx.ensurePullRequest, pipeline.And(repoCtx.hasCommits(), pipeline.Bool(createPR))),
	)
	p.WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instrumentation.PipelineForRepositoryCompleted(r, result.Err())
		return result.Err()
	})
	return p
}

func (c *Command) updateReposInParallel() pipeline.Supplier {
	return func(ctx context.Context, pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		c.instrumentation.BatchPipelineStarted(c.repositories)
		for _, r := range c.GetRepositories() {
			select {
			case <-ctx.Done():
				return
			default:
				p := c.createPipeline(r)
				pipelines <- p
			}
		}
	}
}

func (c *Command) configureInfrastructure(_ context.Context) error {
	c.appService.ConfigureInfrastructure()
	return nil
}

func (c *Command) fetchRepositories(_ context.Context) error {
	repos, err := c.appService.repoStore.FetchGitRepositories()
	c.repositories = repos
	return err
}

func (c *Command) GetRepositories() []*domain.GitRepository {
	return c.repositories
}
