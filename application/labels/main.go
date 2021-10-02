package labels

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to keep repository labels in sync.
	Command struct {
		cfg             *cfg.Configuration
		cliCommand      *cli.Command
		appService      *AppService
		repos           []*domain.GitRepository
		console         logr.Logger
		instrumentation instrumentation.BatchInstrumentation
	}
)

// NewCommand returns a new instance.
func NewCommand(
	cfg *cfg.Configuration,
	appService *AppService,
	instrumentation instrumentation.BatchInstrumentation,
) *Command {
	c := &Command{
		cfg:             cfg,
		appService:      appService,
		instrumentation: instrumentation,
	}
	return c
}

func (c *Command) runCommand(_ *cli.Context) error {
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("configure infrastructure", c.configureInfrastructure),
		pipeline.NewStepFromFunc("fetch repositories", c.fetchRepositories),
		parallel.NewWorkerPoolStep("update labels for all repos", c.cfg.Project.Jobs, c.updateRepos(), c.instrumentation.NewCollectErrorHandler(c.repos, c.cfg.Project.SkipBroken)),
	).Run()
	return result.Err
}

func (c *Command) updateRepos() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		c.instrumentation.BatchPipelineStarted(c.repos)
		for _, r := range c.repos {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	uc := &labelUseCase{
		appService: c.appService,
	}
	return pipeline.NewPipeline().AddBeforeHook(c.appService.factory.NewPipelineLogger("")).WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ pipeline.Context) error {
			c.instrumentation.PipelineForRepositoryStarted(r)
			return nil
		}),
		pipeline.NewStep("fetch labels", uc.fetchLabelsForRepository(r)),
		pipeline.NewStep("determine which labels to modify", c.determineLabelsToModify(uc, r)),
		pipeline.NewStep("determine which labels to delete", c.determineLabelsToDelete(uc, r)),
		pipeline.NewStep("update existing labels", uc.updateLabelsForRepositoryAction(r)),
		pipeline.NewStep("delete unwanted labels", uc.deleteLabelsForRepository(r)),
	).WithFinalizer(func(ctx pipeline.Context, result pipeline.Result) error {
		c.instrumentation.PipelineForRepositoryCompleted(r, result.Err)
		result.Name = r.URL.GetFullName()
		return result.Err
	})
}

func (c *Command) determineLabelsToModify(uc *labelUseCase, r *domain.GitRepository) pipeline.ActionFunc {
	converter := cfg.RepositoryLabelSetConverter{}
	return func(ctx pipeline.Context) pipeline.Result {
		toModify, err := converter.ConvertToEntity(c.cfg.RepositoryLabels.SelectModifications())
		if err != nil {
			return pipeline.Result{Err: err}
		}
		uc.labelsToModify = toModify

		mergedSet := r.Labels.Merge(uc.labelsToModify)
		err = r.SetLabels(mergedSet)
		return pipeline.Result{Err: err}
	}
}

func (c *Command) determineLabelsToDelete(uc *labelUseCase, r *domain.GitRepository) pipeline.ActionFunc {
	converter := cfg.RepositoryLabelSetConverter{}
	return func(ctx pipeline.Context) pipeline.Result {
		toDelete, err := converter.ConvertToEntity(c.cfg.RepositoryLabels.SelectDeletions())
		if err != nil {
			return pipeline.Result{Err: err}
		}
		uc.labelsToDelete = toDelete

		reducedSet := r.Labels.Without(toDelete)
		err = r.SetLabels(reducedSet)
		return pipeline.Result{Err: err}
	}
}

func (c *Command) fetchRepositories(_ pipeline.Context) error {
	repos, err := c.appService.repoStore.FetchGitRepositories()
	c.repos = repos
	return err
}

func (c *Command) configureInfrastructure(_ pipeline.Context) error {
	c.appService.ConfigureInfrastructure()
	return nil
}
