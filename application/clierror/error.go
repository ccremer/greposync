package clierror

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

var (
	// ErrUsage is an error that is caused by an incorrect usage by the user.
	ErrUsage = errors.New("usage error")

	// ErrPipeline is an error that is caused by an action aborting a pipeline.
	ErrPipeline = errors.New("pipeline error")
)

// NewErrorHandler will print the CLI help and exit with code 2 if the given error is a ErrUsage.
func NewErrorHandler(logger logr.Logger) cli.ExitErrHandlerFunc {
	return func(context *cli.Context, err error) {
		if errors.Is(err, ErrUsage) {
			logger.Error(err, "user error")
			cli.ShowCommandHelpAndExit(context, context.Command.Name, 2)
			return
		}
		var cliErr cli.ExitCoder
		if errors.As(err, &cliErr) {
			os.Exit(cliErr.ExitCode())
		}
	}
}

// AsUsageError returns the given error wrapped in a ErrUsage.
func AsUsageError(err error) error {
	return fmt.Errorf("%w: %v", ErrUsage, err)
}

// AsUsageErrorf returns a ErrUsage with the given format
func AsUsageErrorf(format string, a ...interface{}) error {
	return AsUsageError(fmt.Errorf(format, a...))
}

// AsFlagUsageError returns a flag ErrUsage with the given error.
func AsFlagUsageError(flagName string, err error) error {
	return AsUsageError(fmt.Errorf("invalid flag --%s: %v", flagName, err))
}

// AsFlagUsageErrorf returns a flag ErrUsage with the given message.
func AsFlagUsageErrorf(flagName, format string, a ...interface{}) error {
	return AsUsageErrorf("invalid flag --%s: %s", flagName, fmt.Sprintf(format, a...))
}
