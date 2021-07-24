package initialize

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to initialize a new template repository.
	Command struct {
		cfg           *cfg.Configuration
		cliCommand    *cli.Command
		configFiles   map[string][]byte
		templateFiles map[string][]byte
	}
)

// NewCommand returns a new instance.
func NewCommand(cfg *cfg.Configuration) *Command {
	c := &Command{
		cfg: cfg,
		configFiles: map[string][]byte{
			"greposync.yml":       grepoSyncYml,
			"config_defaults.yml": configDefaultsYml,
			"managed_repos.yml":   managedReposYml,
		},
		templateFiles: map[string][]byte{
			cfg.Template.RootDir + "/_helpers.tpl": helperTpl,
			cfg.Template.RootDir + "/README.md":    readmeTpl,
		},
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
		Name:   "init",
		Usage:  "Initializes a template repository in the current working directory",
		Action: c.runCommand,
	}
}

func (c *Command) runCommand(_ *cli.Context) error {
	logger := printer.PipelineLogger{Logger: printer.New().SetLevel(printer.DefaultLevel)}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("create main config files", c.createMainConfigFiles()),
		pipeline.NewStep("create template dir", c.createTemplateDir()),
		pipeline.NewStep("create template files", c.createTemplateFiles()),
	).Run()
	return result.Err
}
