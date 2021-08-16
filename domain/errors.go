package domain

import "errors"

var ErrInvalidArgument = errors.New("invalid argument")

func hasSucceeded(err error) bool {
	return err == nil
}

func hasFailed(err error) bool {
	return err != nil
}

func firstOf(errors ...error) error {
	for _, err := range errors {
		if hasFailed(err) {
			return err
		}
	}
	return nil
}
