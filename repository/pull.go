package repository

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Pull(repo *git.Repository) {

	// Get the working directory for the repository
	w, err := repo.Worktree()
	CheckIfError(err)

	Info("git pull origin")
	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	CheckIfError(err)
	commit, err := repo.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}
