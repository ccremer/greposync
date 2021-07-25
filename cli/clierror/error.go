package clierror

import (
	"errors"
	"fmt"

	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

var (
	// UsageError is an error that is caused by an incorrect usage by the user.
	UsageError = errors.New("usage error")
	// ErrorHandler will print the CLI help and exit with code 2 if the given error is a UsageError.
	ErrorHandler cli.ExitErrHandlerFunc = func(context *cli.Context, err error) {
		if errors.Is(err, UsageError) {
			printer.DefaultPrinter.ErrorF("%v\n", err.Error())
			cli.ShowCommandHelpAndExit(context, context.Command.Name, 2)
			return
		}
	}
)

// AsUsageError returns the given error wrapped in a UsageError.
func AsUsageError(err error) error {
	return fmt.Errorf("%w: %v", UsageError, err)
}

// AsUsageErrorf returns a UsageError with the given format
func AsUsageErrorf(format string, a ...interface{}) error {
	return AsUsageError(fmt.Errorf(format, a...))
}

// AsFlagUsageError returns a flag UsageError with the given error.
func AsFlagUsageError(flagName string, err error) error {
	return AsUsageError(fmt.Errorf("invalid flag --%s: %v", flagName, err))
}

// AsFlagUsageErrorf returns a flag UsageError with the given message.
func AsFlagUsageErrorf(flagName, format string, a ...interface{}) error {
	return AsUsageErrorf("invalid flag --%s: %s", flagName, fmt.Sprintf(format, a...))
}
