package labels

import (
	"context"

	pipeline "github.com/ccremer/go-command-pipeline"
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
		pipeline.NewWorkerPoolStep("update labels for all repos", c.cfg.Project.Jobs, c.updateRepos(), c.instrumentation.NewCollectErrorHandler(c.cfg.Project.SkipBroken)),
	).Run()
	return result.Err()
}

func (c *Command) updateRepos() pipeline.Supplier {
	return func(ctx context.Context, pipelinesCH chan *pipeline.Pipeline) {
		defer close(pipelinesCH)
		c.instrumentation.BatchPipelineStarted(c.repos)
		for _, r := range c.repos {
			select {
			case <-ctx.Done():
				return
			default:
				p := c.createPipeline(r)
				pipelinesCH <- p
			}
		}
	}
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	uc := &labelPipeline{
		appService: c.appService,
		repo:       r,
	}
	return pipeline.NewPipeline().AddBeforeHook(c.appService.factory.NewPipelineLogger("").Accept).WithSteps(
		pipeline.NewStepFromFunc("setup instrumentation", func(_ context.Context) error {
			c.instrumentation.PipelineForRepositoryStarted(r)
			return nil
		}),
		pipeline.NewStepFromFunc("fetch labels", uc.fetchLabelsForRepository),
		pipeline.NewStepFromFunc("determine which labels to modify", c.determineLabelsToModify(uc, r)),
		pipeline.NewStepFromFunc("determine which labels to delete", c.determineLabelsToDelete(uc, r)),
		pipeline.NewStepFromFunc("update existing labels", uc.updateLabelsForRepository),
		pipeline.NewStepFromFunc("delete unwanted labels", uc.deleteLabelsForRepository),
	).WithFinalizer(func(ctx context.Context, result pipeline.Result) error {
		c.instrumentation.PipelineForRepositoryCompleted(r, result.Err())
		return result.Err()
	})
}

func (c *Command) determineLabelsToModify(uc *labelPipeline, r *domain.GitRepository) func(ctx context.Context) error {
	converter := cfg.RepositoryLabelSetConverter{}
	return func(ctx context.Context) error {
		toModify, err := converter.ConvertToEntity(c.cfg.RepositoryLabels.SelectModifications())
		if err != nil {
			return err
		}
		uc.labelsToModify = toModify

		mergedSet := r.Labels.Merge(uc.labelsToModify)
		err = r.SetLabels(mergedSet)
		return err
	}
}

func (c *Command) determineLabelsToDelete(uc *labelPipeline, r *domain.GitRepository) func(ctx context.Context) error {
	converter := cfg.RepositoryLabelSetConverter{}
	return func(ctx context.Context) error {
		toDelete, err := converter.ConvertToEntity(c.cfg.RepositoryLabels.SelectDeletions())
		if err != nil {
			return err
		}
		uc.labelsToDelete = toDelete

		reducedSet := r.Labels.Without(toDelete)
		err = r.SetLabels(reducedSet)
		return nil
	}
}

func (c *Command) fetchRepositories(_ context.Context) error {
	repos, err := c.appService.repoStore.FetchGitRepositories()
	c.repos = repos
	return err
}

func (c *Command) configureInfrastructure(_ context.Context) error {
	c.appService.ConfigureInfrastructure()
	return nil
}

func (c *Command) GetRepositories() []*domain.GitRepository {
	return c.repos
}
