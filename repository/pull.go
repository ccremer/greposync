package repository

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func (s *Service) Pull() {

	// Get the working directory for the repository
	w, err := s.r.Worktree()
	CheckIfError(err)

	Info("git pull origin")
	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	// Print the latest commit that was just pulled
	ref, err := s.r.Head()
	CheckIfError(err)
	commit, err := s.r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}
