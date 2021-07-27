package labels

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/flags"
	"github.com/ccremer/greposync/core/labels"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to keep repository labels in sync.
	Command struct {
		cfg        *cfg.Configuration
		cliCommand *cli.Command
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
		Flags: flags.CombineWithGlobalFlags(
			flags.NewProjectIncludeFlag(),
			flags.NewProjectExcludeFlag(),
		),
	}
}

func (c *Command) runCommand(ctx *cli.Context) error {
	labelService := labels.NewService(repository.NewGitRepositoryProvider(c.cfg))
	return labelService.RunPipeline()
}
