package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to keep repository labels in sync.
	Command struct {
		cfg        *cfg.Configuration
		cliCommand *cli.Command
		appService *AppService
		repos      []*domain.GitRepository
	}
)

// NewCommand returns a new instance.
func NewCommand(cfg *cfg.Configuration, appService *AppService) *Command {
	c := &Command{
		cfg:        cfg,
		appService: appService,
	}
	c.cliCommand = c.createCommand()
	return c
}

func (c *Command) runCommand(_ *cli.Context) error {
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("configure infrastructure", c.configureInfrastructure),
		pipeline.NewStepFromFunc("fetch repositories", c.fetchRepositories),
		parallel.NewWorkerPoolStep("update labels for all repos", c.cfg.Project.Jobs, c.updateRepos(), c.errorHandler()),
	).Run()
	return result.Err
}

func (c *Command) updateRepos() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repos {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) errorHandler() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, repo := range c.repos {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", repo.URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	uc := &labelUseCase{
		appService: c.appService,
	}
	return pipeline.NewPipeline().WithSteps(
		pipeline.NewStep("fetch labels", uc.fetchLabelsForRepository(r)),
		pipeline.NewStep("determine which labels to modify", c.determineLabelsToModify(uc, r)),
		pipeline.NewStep("determine which labels to delete", c.determineLabelsToDelete(uc, r)),
		pipeline.NewStep("update existing labels", uc.updateLabelsForRepositoryAction(r)),
		pipeline.NewStep("delete unwanted labels", uc.deleteLabelsForRepository(r)),
	)
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
