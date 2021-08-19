package domain

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrInvalidArgument is an error that indicates that a particular field is invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrKeyNotFound is an error that indicates that a particular key was not found.
var ErrKeyNotFound = errors.New("key not found")

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

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func checkIfArgumentNil(v interface{}, fieldName string) error {
	if isNil(v) {
		return fmt.Errorf("%w: %s cannot be nil", ErrInvalidArgument, fieldName)
	}
	return nil
}
