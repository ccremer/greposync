package labels

import (
	"github.com/ccremer/greposync/application/flags"
	"github.com/urfave/cli/v2"
)

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
