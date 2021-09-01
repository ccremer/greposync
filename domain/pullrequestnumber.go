package domain

import "fmt"

type PullRequestNumber int

func NewPullRequestNumber(nr *int) *PullRequestNumber {
	if nr == nil {
		return nil
	}
	pnr := PullRequestNumber(*nr)
	return &pnr
}

func (nr PullRequestNumber) String() string {
	return fmt.Sprintf("#%d", nr)
}

func (nr *PullRequestNumber) Int() *int {
	if nr == nil {
		return nil
	}
	v := int(*nr)
	return &v
}
