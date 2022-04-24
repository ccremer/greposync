package test

import (
	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/urfave/cli/v2"
)

type Command struct {
	cfg          *cfg.Configuration
	repositories []*domain.GitRepository
	appService   *AppService
	instr        instrumentation.BatchInstrumentation
	logFactory   logging.LoggerFactory

	exitOnFail bool
}

// NewCommand returns a new Command instance.
func NewCommand(
	cfg *cfg.Configuration,
	configurator *AppService,
	factory logging.LoggerFactory,
	instrumentation instrumentation.BatchInstrumentation,
) *Command {
	c := &Command{
		cfg:        cfg,
		appService: configurator,
		instr:      instrumentation,
		logFactory: factory,
	}
	return c
}

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.createCliCommand()
}

func (c *Command) createCliCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Test the rendered template against test cases",
		Description: `Test cases are defined as local, simulated repositories in 'tests' directory, where each subdirectory itself is a separate test case.
The expected file structure should resemble this format: 

tests
└── case-1 
    ├── <file>
    └── .sync.yml

'case-1' is the test case name.
<file> represents any files that are to be rendered (for example README.md, Makefile etc.) with their contents being the desired output.
'.sync.yml' is the sync config for this simulated repository and it works exactly as the .sync.yml syntax in real repositories.

When running this subcommand, these test cases are picked up and its template output rendered in a new directory '.tests'.
A 'git diff' will be computed and if it's non-empty, the test case is considered failed.

This command can be used to verify that the template is correct before rolling it out to production repositories.
`,
		Action: c.runCommand,
		Before: c.validateTestCommand,
		Flags: flags.CombineWithGlobalFlags(
			&cli.BoolFlag{
				Name:        "exit-code",
				Destination: &c.exitOnFail,
				Usage:       "Exits app with exit code 3 if a test case failed.",
			},
		),
	}
}
