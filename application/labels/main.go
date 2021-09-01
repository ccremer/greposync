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
		pipeline.NewStep("update existing labels", uc.updateLabelsForRepository(r)),
		pipeline.NewStep("delete unwanted labels", uc.deleteLabelsForRepository(r)),
	)
}

func (c *Command) fetchRepositories(_ pipeline.Context) error {
	repos, err := c.appService.repoStore.FetchGitRepositories()
	c.repos = repos
	return err
}
