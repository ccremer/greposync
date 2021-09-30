package domain

import "fmt"

// PullRequestNumber identifies a PullRequest by a number in a Git hosting service.
type PullRequestNumber int

// NewPullRequestNumber takes the given number and returns a new instance.
// If nr is nil, then nil is returned.
func NewPullRequestNumber(nr *int) *PullRequestNumber {
	if nr == nil {
		return nil
	}
	pnr := PullRequestNumber(*nr)
	return &pnr
}

// String returns the number prefixed with `#`.
func (nr PullRequestNumber) String() string {
	return fmt.Sprintf("#%d", nr)
}

// Int returns nil if nr is also nil.
// Otherwise, it returns an int pointer.
func (nr *PullRequestNumber) Int() *int {
	if nr == nil {
		return nil
	}
	v := int(*nr)
	return &v
}
