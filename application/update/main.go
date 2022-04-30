package update

import (
	"context"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/urfave/cli/v2"
)

const (
	dryRunFlagName    = "dry-run"
	prCreateFlagName  = "pr.create"
	prBodyFlagName    = "pr.bodyTemplate"
	amendFlagName     = "git.amend"
	forcePushFlagName = "git.forcePush"
	showDiffFlagName  = "log.showDiff"
)

type (
	// Command is a facade service for the update command that holds all dependent services and settings.
	Command struct {
		cfg          *cfg.Configuration
		repositories []*domain.GitRepository
		appService   *AppService
		instr        instrumentation.BatchInstrumentation
		logFactory   logging.LoggerFactory

		dryRunFlag string
		PrLabels   cli.StringSlice
	}
)

// NewCommand returns a new Command instance.
func NewCommand(
	cfg *cfg.Configuration,
	configurator *AppService,
	factory logging.LoggerFactory,
	instrumentation instrumentation.BatchInstrumentation,
) *Command {
	c := &Command{
		cfg:        cfg,
		appService: configurator,
		instr:      instrumentation,
		logFactory: factory,
	}
	return c
}

func (c *Command) runCommand(cliCtx *cli.Context) error {
	logger := c.logFactory.NewPipelineLogger("")
	ctx := pipeline.MutableContext(cliCtx.Context)
	p := pipeline.NewPipeline().AddBeforeHook(logger.Accept).WithSteps(
		pipeline.NewStepFromFunc("configure infrastructure", c.configureInfrastructure),
		pipeline.NewStepFromFunc("fetch managed repos config", c.fetchRepositories),
		pipeline.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.instr.NewCollectErrorHandler(c.cfg.Project.SkipBroken)),
	)
	p.WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instr.BatchPipelineCompleted("Update finished", c.repositories)
		return result.Err()
	})
	return p.RunWithContext(ctx).Err()
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {

	resetRepo := !c.cfg.Git.SkipReset
	enabledCommits := !c.cfg.Git.SkipCommit
	enabledPush := !c.cfg.Git.SkipPush
	showDiff := c.cfg.Log.ShowDiff
	createPR := c.cfg.PullRequest.Create

	up := &updatePipeline{
		repo:       r,
		appService: c.appService,
		prLabels:   c.PrLabels.Value(),
	}

	logger := c.logFactory.NewPipelineLogger(r.URL.GetFullName())
	pipe := pipeline.NewPipeline().AddBeforeHook(logger.Accept)
	pipe.WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ context.Context) error {
			c.instr.PipelineForRepositoryStarted(up.repo)
			return nil
		}),

		pipeline.NewPipeline().AddBeforeHook(logger.Accept).
			WithNestedSteps("prepare workspace",

				pipeline.ToStep("clone repository", up.clone, up.dirMissing()),
				pipeline.ToStep("fetch", up.fetch, pipeline.Bool(resetRepo)),
				pipeline.ToStep("reset", up.reset, pipeline.Bool(resetRepo)),
				pipeline.ToStep("checkout branch", up.checkout, pipeline.Bool(resetRepo)),
				pipeline.ToStep("pull", up.pull, pipeline.Bool(resetRepo)),
			),

		pipeline.NewPipeline().AddBeforeHook(logger.Accept).
			WithNestedSteps("render",
				pipeline.NewStepFromFunc("render templates", up.renderTemplates),
				pipeline.NewStepFromFunc("cleanup unwanted files", up.cleanupUnwantedFiles),
			),

		pipeline.If(pipeline.And(pipeline.Bool(enabledCommits), up.isDirty()),
			pipeline.NewPipeline().
				AddBeforeHook(logger.Accept).
				WithNestedSteps("commit changes",
					pipeline.NewStepFromFunc("add", up.add),
					pipeline.NewStepFromFunc("commit", up.commit),
				),
		),

		pipeline.ToStep("show diff", up.diff, pipeline.Bool(showDiff)),
		pipeline.ToStep("push changes", up.push, pipeline.And(pipeline.Bool(enabledPush), up.hasCommits())),
		pipeline.ToStep("find existing pull request", up.fetchPullRequest, pipeline.Bool(createPR)),
		pipeline.ToStep("ensure pull request", up.ensurePullRequest, pipeline.And(up.hasCommits(), pipeline.Bool(createPR))),
	)
	pipe.WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instr.PipelineForRepositoryCompleted(r, result.Err())
		return result.Err()
	})
	return pipe
}

func (c *Command) updateReposInParallel() pipeline.Supplier {
	return func(ctx context.Context, pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		c.instr.BatchPipelineStarted("Update started", c.repositories)
		for _, r := range c.repositories {
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

func (c *Command) fetchRepositories(ctx context.Context) error {
	repos, err := c.appService.repoStore.FetchGitRepositories()
	c.repositories = repos
	pipeline.StoreInContext(ctx, instrumentation.RepositoriesContextKey{}, repos)
	return err
}
