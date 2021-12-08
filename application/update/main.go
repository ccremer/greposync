package update

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
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
	p := pipeline.NewPipeline().WithContext(c).AddBeforeHook(logger).WithSteps(
		pipeline.NewStep("configure infrastructure", c.configureInfrastructure()),
		pipeline.NewStep("fetch managed repos config", c.fetchRepositories()),
		parallel.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.instrumentation.NewCollectErrorHandler(c.cfg.Project.SkipBroken)),
	)
	p.WithFinalizer(func(ctx pipeline.Context, result pipeline.Result) error {
		c.instrumentation.BatchPipelineCompleted(c.GetRepositories())
		return result.Err
	})
	return p.Run().Err
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {

	resetRepo := !c.cfg.Git.SkipReset
	enabledCommits := !c.cfg.Git.SkipCommit
	enabledPush := !c.cfg.Git.SkipPush
	showDiff := c.cfg.Log.ShowDiff
	createPR := c.cfg.PullRequest.Create

	repoCtx := &pipelineContext{
		repo:       r,
		appService: c.appService,
	}

	logger := c.logFactory.NewPipelineLogger(r.URL.GetFullName())
	p := pipeline.NewPipeline().AddBeforeHook(logger)
	p.WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ pipeline.Context) error {
			c.instrumentation.PipelineForRepositoryStarted(repoCtx.repo)
			return nil
		}),

		pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("prepare workspace",
				predicate.ToStep("clone repository", repoCtx.clone(), repoCtx.dirMissing()),
				predicate.ToStep("fetch", repoCtx.fetch(), predicate.Bool(resetRepo)),
				predicate.ToStep("reset", repoCtx.reset(), predicate.Bool(resetRepo)),
				predicate.ToStep("checkout branch", repoCtx.checkout(), predicate.Bool(resetRepo)),
				predicate.ToStep("pull", repoCtx.pull(), predicate.Bool(resetRepo)),
			),

		pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("render",
				pipeline.NewStep("render templates", repoCtx.renderTemplates()),
				pipeline.NewStep("cleanup unwanted files", repoCtx.cleanupUnwantedFiles()),
			),

		predicate.WrapIn(pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("commit changes",
				pipeline.NewStep("add", repoCtx.add()),
				pipeline.NewStep("commit", repoCtx.commit()),
			),
			predicate.And(predicate.Bool(enabledCommits), repoCtx.isDirty())),

		predicate.ToStep("show diff", repoCtx.diff(), predicate.Bool(showDiff)),
		predicate.ToStep("push changes", repoCtx.push(), predicate.And(predicate.Bool(enabledPush), repoCtx.hasCommits())),
		predicate.ToStep("find existing pull request", repoCtx.fetchPullRequest(), predicate.Bool(createPR)),
		predicate.ToStep("ensure pull request", repoCtx.ensurePullRequest(), predicate.And(repoCtx.hasCommits(), predicate.Bool(createPR))),
	)
	p.WithFinalizer(func(ctx pipeline.Context, result pipeline.Result) error {
		c.instrumentation.PipelineForRepositoryCompleted(r, result.Err)
		result.Name = r.URL.GetFullName()
		return result.Err
	})
	return p
}

func (c *Command) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		c.instrumentation.BatchPipelineStarted(c.repositories)
		for _, r := range c.GetRepositories() {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) configureInfrastructure() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		c.appService.ConfigureInfrastructure()
		return pipeline.Result{}
	}
}

func (c *Command) fetchRepositories() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		repos, err := c.appService.repoStore.FetchGitRepositories()
		c.repositories = repos
		return pipeline.Result{Err: err}
	}
}

func (c *Command) GetRepositories() []*domain.GitRepository {
	return c.repositories
}
