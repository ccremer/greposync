package repository

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func MakeCommit(repo *git.Repository) {
	w, err := repo.Worktree()
	CheckIfError(err)

	Info("git add *")
	err = w.AddWithOptions(&git.AddOptions{All: true})
	CheckIfError(err)

	Info("git status --porcelain")
	fmt.Println(w.Status())

	Info("git commit -m \"New update from template\"")
	commit, err := w.Commit("New update from template", &git.CommitOptions{
		Author: &object.Signature{
			Name: "just me",
			Email: "john@doe.org",
			When: time.Now(),
		},
	})

	obj, err := repo.CommitObject(commit)
	fmt.Println(obj)
	CheckIfError(err)
}
