package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/flags"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/labels"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to keep repository labels in sync.
	Command struct {
		cfg        *cfg.Configuration
		cliCommand *cli.Command
		repoStore  core.GitRepositoryStore
	}
)

// NewCommand returns a new instance.
func NewCommand(cfg *cfg.Configuration, repoStore core.GitRepositoryStore) *Command {
	c := &Command{
		cfg:       cfg,
		repoStore: repoStore,
	}
	c.cliCommand = c.createCommand()
	return c
}

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.cliCommand
}

func (c *Command) createCommand() *cli.Command {
	return &cli.Command{
		Name:   "labels",
		Usage:  "Synchronizes repository labels",
		Before: c.validateCommand,
		Action: c.runCommand,
		Flags:  flags.CombineWithGlobalFlags(),
	}
}

func (c *Command) runCommand(_ *cli.Context) error {
	repos, err := c.repoStore.FetchGitRepositories()
	if err != nil {
		return err
	}
	result := pipeline.NewPipeline().WithSteps(
		parallel.NewWorkerPoolStep("update labels for all repos", c.cfg.Project.Jobs, c.updateRepos(repos), c.errorHandler(repos)),
	).Run()
	return result.Err
}

func (c *Command) updateRepos(repos []core.GitRepository) parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range repos {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) errorHandler(repos []core.GitRepository) parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, service := range repos {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", service.GetConfig().URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}

func (c *Command) updateLabelsForRepo(url *core.GitURL) pipeline.ActionFunc {
	return c.fireEvent(url, labels.LabelUpdateEvent)
}

func (c *Command) fireEvent(u *core.GitURL, event core.EventName) pipeline.ActionFunc {
	return func() pipeline.Result {
		result := <-core.FireEvent(event, core.EventSource{
			Url: u,
		})
		return pipeline.Result{Err: result.Error}
	}
}

func (c *Command) createPipeline(r core.GitRepository) *pipeline.Pipeline {
	return pipeline.NewPipeline().WithSteps(
		pipeline.NewStep("update labels", c.updateLabelsForRepo(r.GetConfig().URL)),
	)
}
