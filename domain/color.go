package domain

import (
	"fmt"
	"regexp"
)

// Color is a 6-digit uppercase hexadecimal string value with '#' prefix
type Color string

var colorRegex *regexp.Regexp

func init() {
	colorRegex = regexp.MustCompile("^#[A-F0-9]{6}$")
}

func (c Color) String() string {
	return string(c)
}

func (c Color) CheckValue() error {
	if colorRegex.MatchString(c.String()) {
		return nil
	}
	return fmt.Errorf("%w: color value must be 6-digit uppercase hexadecimal with '#' prefix: %s", ErrInvalidArgument, c)
}
