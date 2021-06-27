package repository

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (s *Service) SwitchBranch() {
	s.CheckoutBranch()
}

func (s *Service) CheckoutBranch() {
	branch := fmt.Sprintf("refs/heads/%s", "my-branch")
	b := plumbing.ReferenceName(branch)

	w, err := s.r.Worktree()
	CheckIfError(err)

	// First try to checkout branch
	err = w.Checkout(&git.CheckoutOptions{Create: false, Force: true, Branch: b})
	if err != nil {
		// got an error  - try to create it
		err = w.Checkout(&git.CheckoutOptions{Create: true, Force: true, Branch: b})
		CheckIfError(err)
	}
}
