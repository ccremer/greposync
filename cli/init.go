package cli

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/initialize"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

type (
	// InitCommand contains the logic to initialize a new template repository.
	InitCommand struct {
		cfg        *cfg.Configuration
		cliCommand *cli.Command
	}
)

// NewInitCommand returns a new instance.
func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

func (c *InitCommand) createInitCommand() *cli.Command {
	c.cliCommand = &cli.Command{
		Name:   "init",
		Usage:  "Initializes a template repository in the current working directory",
		Action: c.runInitCommand,
	}
	return c.cliCommand
}

func (c *InitCommand) runInitCommand(_ *cli.Context) error {
	logger := printer.PipelineLogger{Logger: printer.New().SetLevel(printer.DefaultLevel)}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("create main config files", initialize.CreateMainConfigFiles()),
		pipeline.NewStep("create template dir", initialize.CreateTemplateDir()),
		pipeline.NewStep("create template files", initialize.CreateTemplateFiles()),
	).Run()
	return result.Err
}
