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
		Name:   "test",
		Usage:  "Test the rendered template against test cases",
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
