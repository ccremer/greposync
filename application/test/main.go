package test

import (
	"context"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/domain"
	"github.com/urfave/cli/v2"
)

func (c *Command) runCommand(cliCtx *cli.Context) error {
	logger := c.logFactory.NewPipelineLogger("")
	ctx := pipeline.MutableContext(cliCtx.Context)
	p := pipeline.NewPipeline().AddBeforeHook(logger.Accept).WithSteps(
		pipeline.NewStepFromFunc("configure infrastructure", c.configureInfrastructure),
		pipeline.NewStepFromFunc("fetch test repositories", c.fetchRepositories),
		pipeline.NewStepFromFunc("clean test directory", c.cleanTestDir),
		pipeline.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.instr.NewCollectErrorHandler(true)),
	)
	p.WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instr.BatchPipelineCompleted("Tests finished", c.repositories)
		return result.Err()
	})
	return p.RunWithContext(ctx).Err()
}

func (c *Command) createRepositoryPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	up := &updatePipeline{
		repo:               r,
		appService:         c.appService,
		failPipelineIfDiff: c.exitOnFail,
	}

	logger := c.logFactory.NewPipelineLogger(r.URL.GetFullName())
	pipe := up.AddBeforeHook(logger.Accept)
	pipe.WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ context.Context) error {
			c.instr.PipelineForRepositoryStarted(up.repo)
			return nil
		}),
		pipeline.NewStepFromFunc("create dir", up.createOutputDir),
		pipeline.NewStepFromFunc("copy sync file", up.copySyncFile),
		pipeline.NewPipeline().AddBeforeHook(logger.Accept).
			WithNestedSteps("render",
				pipeline.NewStepFromFunc("render templates", up.renderTemplates),
				pipeline.NewStepFromFunc("cleanup unwanted files", up.cleanupUnwantedFiles),
			),

		pipeline.NewStepFromFunc("show diff", up.diff),
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
		c.instr.BatchPipelineStarted("Tests started", c.repositories)
		for _, r := range c.repositories {
			select {
			case <-ctx.Done():
				return
			default:
				p := c.createRepositoryPipeline(r)
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

func (c *Command) cleanTestDir(_ context.Context) error {
	return os.RemoveAll(c.appService.repoStore.TestOutputRootDir)
}
