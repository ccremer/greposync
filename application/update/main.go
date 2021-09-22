package update

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

const (
	dryRunFlagName   = "dry-run"
	prCreateFlagName = "pr-create"
	prBodyFlagName   = "pr-bodyTemplate"
	amendFlagName    = "git-amend"
	showDiffFlagName = "log-showDiff"
)

type (
	// Command is a facade service for the update command that holds all dependent services and settings.
	Command struct {
		cfg             *cfg.Configuration
		repositories    []*domain.GitRepository
		appService      *AppService
		instrumentation *UpdateInstrumentation
		logFactory      logging.LoggerFactory
	}
)

// NewCommand returns a new Command instance.
func NewCommand(
	cfg *cfg.Configuration,
	configurator *AppService,
	factory logging.LoggerFactory,
	instrumentation *UpdateInstrumentation,
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
	p := pipeline.NewPipeline().AddBeforeHook(logger).WithSteps(
		pipeline.NewStep("configure infrastructure", c.configureInfrastructure()),
		pipeline.NewStep("fetch managed repos config", c.fetchRepositories()),
		parallel.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.collectErrors()),
	)
	p.WithFinalizer(func(result pipeline.Result) error {
		c.instrumentation.batchPipelineCompleted(c.repositories)
		return result.Err
	})
	return p.Run().Err
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	sc := &cfg.SyncConfig{
		PullRequest: c.cfg.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: c.cfg.Template.RootDir,
		},
	}

	resetRepo := !c.cfg.Git.SkipReset
	enabledCommits := !c.cfg.Git.SkipCommit
	enabledPush := !c.cfg.Git.SkipPush
	showDiff := c.cfg.Log.ShowDiff

	repoCtx := &pipelineContext{
		repo:       r,
		appService: c.appService,
	}

	logger := c.logFactory.NewPipelineLogger(r.URL.GetFullName())
	p := pipeline.NewPipeline().AddBeforeHook(logger)
	p.WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ pipeline.Context) error {
			c.instrumentation.pipelineForRepositoryStarted(repoCtx.repo)
			return nil
		}),
		pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("prepare workspace",
				predicate.ToStep("clone repository", repoCtx.clone(), repoCtx.dirMissing()),
				predicate.ToStep("fetch", repoCtx.fetch(), predicate.Bool(resetRepo)),
				predicate.ToStep("reset", repoCtx.reset(), predicate.Bool(resetRepo)),
				predicate.ToStep("checkout branch", repoCtx.checkout(), predicate.Bool(resetRepo)),
				predicate.ToStep("pull", repoCtx.fetch(), predicate.Bool(resetRepo)),
			),
		pipeline.NewStep("render templates", repoCtx.renderTemplates()),
		predicate.WrapIn(pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("commit changes",
				pipeline.NewStep("add", repoCtx.add()),
				pipeline.NewStep("commit", repoCtx.commit()),
				predicate.ToStep("show diff", repoCtx.diff(), predicate.Bool(showDiff)),
			),
			predicate.And(predicate.Bool(enabledCommits), repoCtx.isDirty())),
		predicate.ToStep("push changes", repoCtx.push(), predicate.And(predicate.Bool(enabledPush), repoCtx.hasCommits())),
		predicate.ToStep("find existing pull request", repoCtx.fetchPullRequest(), predicate.Bool(sc.PullRequest.Create)),
		predicate.ToStep("update pull request", repoCtx.ensurePullRequest(), predicate.And(repoCtx.hasCommits(), predicate.Bool(sc.PullRequest.Create))),
	)
	p.WithFinalizer(func(result pipeline.Result) error {
		c.instrumentation.pipelineForRepositoryCompleted(r, result.Err)
		result.Name = repoCtx.repo.URL.GetFullName()
		return result.Err
	})
	return p
}

func (c *Command) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		c.instrumentation.batchPipelineStarted(len(c.repositories))
		for _, r := range c.repositories {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) collectErrors() parallel.ResultHandler {
	if c.cfg.Project.SkipBroken {
		return c.ignoreErrors()
	}
	return c.reduceErrors()
}

func (c *Command) ignoreErrors() parallel.ResultHandler {
	// Do not propagate update errors from single repositories up the stack
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		c.instrumentation.results = results
		return pipeline.Result{}
	}
}

func (c *Command) reduceErrors() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		c.instrumentation.results = results
		var err error
		for index, repo := range c.repositories {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", repo.URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: fmt.Errorf("%w: %s", clierror.ErrPipeline, err)}
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
