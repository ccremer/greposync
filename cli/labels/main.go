package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/flags"
	"github.com/ccremer/greposync/printer"
	"github.com/ccremer/greposync/repository"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to keep repository labels in sync.
	Command struct {
		cfg          *cfg.Configuration
		cliCommand   *cli.Command
		repoServices []*repository.Service
	}
)

// NewCommand returns a new instance.
func NewCommand(cfg *cfg.Configuration) *Command {
	c := &Command{
		cfg: cfg,
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
		Flags:  flags.CombineWithGlobalFlags(
		//projectIncludeFlag,
		//projectExcludeFlag,
		),
	}
}

func (c *Command) runCommand(ctx *cli.Context) error {
	logger := printer.PipelineLogger{Logger: printer.New().SetName(ctx.Command.Name).SetLevel(printer.DefaultLevel)}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("parse config", c.parseServices()),
		parallel.NewWorkerPoolStep("update labels", c.cfg.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	).Run()
	return result.Err
}

func (c *Command) parseServices() func() pipeline.Result {
	return func() pipeline.Result {
		s, err := repository.NewServicesFromFile(c.cfg)
		c.repoServices = s
		return pipeline.Result{Err: err}
	}
}

func (c *Command) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repoServices {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) createPipeline(r *repository.Service) *pipeline.Pipeline {

	log := printer.New().SetName(r.Config.Name).SetLevel(printer.DefaultLevel)
	logger := printer.PipelineLogger{Logger: log}

	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewStep("prepare API", r.InitializeGitHubProvider(c.cfg.PullRequest)),
		pipeline.NewStep("update labels", r.CreateOrUpdateLabels(c.cfg.RepositoryLabels)),
	)
	return p
}

func (c *Command) errorHandler() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, service := range c.repoServices {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", service.Config.Name, result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}
