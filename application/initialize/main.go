package initialize

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/valuestore"
	"github.com/urfave/cli/v2"
)

type (
	// Command contains the logic to initialize a new template repository.
	Command struct {
		cfg           *cfg.Configuration
		cliCommand    *cli.Command
		configFiles   map[string][]byte
		templateFiles map[string][]byte
		plog          *logging.PipelineLogger

		TemplateDir string
	}
)

// NewCommand returns a new instance.
func NewCommand(cfg *cfg.Configuration, factory logging.LoggerFactory) *Command {
	c := &Command{
		cfg:  cfg,
		plog: factory.NewPipelineLogger("init"),
		configFiles: map[string][]byte{
			"greposync.yml":                 grepoSyncYml,
			valuestore.GlobalConfigFileName: configDefaultsYml,
			"managed_repos.yml":             managedReposYml,
		},
		templateFiles: map[string][]byte{
			cfg.Template.RootDir + "/_helpers.tpl": helperTpl,
			cfg.Template.RootDir + "/README.md":    readmeTpl,
		},
	}
	return c
}

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.createCommand()
}

func (c *Command) createCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initializes a template repository in the current working directory",
		Description: `This command creates '` + valuestore.GlobalConfigFileName + `', 'managed_repositories.yml' and some template files with example content.
If any file already exists, it will be left untouched.`,
		Action: c.runCommand,
		Flags: []cli.Flag{
			flags.NewTemplateRootDirFlag(&c.TemplateDir),
		},
	}
}

func (c *Command) runCommand(_ *cli.Context) error {
	result := pipeline.NewPipeline().AddBeforeHook(c.plog.Accept).WithSteps(
		pipeline.NewStepFromFunc("create main config files", c.createMainConfigFiles),
		pipeline.NewStepFromFunc("create template dir", c.createTemplateDir),
		pipeline.NewStepFromFunc("create template files", c.createTemplateFiles),
	).Run()
	return result.Err()
}
