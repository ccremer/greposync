package cli

import (
	"encoding/json"
	"fmt"
	"regexp"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/ccremer/greposync/repository"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

type (
	// LabelsCommand contains the logic to keep repository labels in sync.
	LabelsCommand struct {
		cfg          *cfg.Configuration
		cliCommand   *cli.Command
		repoServices []*repository.Service
	}
)

// NewLabelsCommand returns a new instance.
func NewLabelsCommand(cfg *cfg.Configuration) *LabelsCommand {
	return &LabelsCommand{
		cfg: cfg,
	}
}

func (c *LabelsCommand) createCommand() *cli.Command {
	c.cliCommand = &cli.Command{
		Name:   "labels",
		Usage:  "Synchronizes repository labels",
		Before: c.validateCommand,
		Action: c.runCommand,
		Flags:  combineWithGlobalFlags(
		//projectIncludeFlag,
		//projectExcludeFlag,
		),
	}
	return c.cliCommand
}

func (c *LabelsCommand) validateCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig(GrepoSyncFileName, config, ctx); err != nil {
		return err
	}

	if err := validateGlobalFlags(ctx); err != nil {
		return err
	}

	if _, err := regexp.Compile(config.Project.Include); err != nil {
		return fmt.Errorf("invalid flag --%s: %v", projectIncludeFlagName, err)
	}
	if _, err := regexp.Compile(config.Project.Exclude); err != nil {
		return fmt.Errorf("invalid flag --%s: %v", projectExcludeFlagName, err)
	}

	if jobs := config.Project.Jobs; jobs > JobsMaximumCount || jobs < JobsMinimumCount {
		return fmt.Errorf("--%s is required to be between %d and %d", projectJobsFlagName, JobsMinimumCount, JobsMaximumCount)
	}

	for key, label := range config.RepositoryLabels {
		if label.Name == "" {
			return fmt.Errorf("label name with key '%s' cannot be empty in '%s'", key, "repositoryLabels")
		}
	}

	config.Sanitize()
	j, _ := json.Marshal(config)
	printer.DebugF("Using config: %s", j)
	return nil
}

func (c *LabelsCommand) runCommand(ctx *cli.Context) error {
	logger := printer.PipelineLogger{Logger: printer.New().SetName(ctx.Command.Name).SetLevel(printer.DefaultLevel)}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("parse config", c.parseServices()),
		parallel.NewWorkerPoolStep("update labels", config.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	).Run()
	return result.Err
}

func (c *LabelsCommand) parseServices() func() pipeline.Result {
	return func() pipeline.Result {
		s, err := repository.NewServicesFromFile(config)
		c.repoServices = s
		return pipeline.Result{Err: err}
	}
}

func (c *LabelsCommand) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repoServices {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *LabelsCommand) createPipeline(r *repository.Service) *pipeline.Pipeline {

	log := printer.New().SetName(r.Config.Name).SetLevel(printer.DefaultLevel)
	logger := printer.PipelineLogger{Logger: log}

	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewStep("prepare API", r.InitializeGitHubProvider(c.cfg.PullRequest)),
		pipeline.NewStep("update labels", r.CreateOrUpdateLabels(c.cfg.RepositoryLabels)),
	)
	return p
}

func (c *LabelsCommand) errorHandler() parallel.ResultHandler {
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
