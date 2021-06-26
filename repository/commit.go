package repository

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func MakeCommit(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_ = w.AddGlob("*")
	fmt.Println(w.Status())

	commit, err := w.Commit("New update from template", &git.CommitOptions{
		Author: &object.Signature{
			Name: "just me",
			Email: "john@doe.org",
			When: time.Now(),
		},
	})

	obj, err := repo.CommitObject(commit)
	fmt.Println(obj)
	return err
}
