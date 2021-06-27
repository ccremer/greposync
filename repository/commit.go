package repository

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func (s *Service) MakeCommit() {
	w, err := s.r.Worktree()
	CheckIfError(err)

	Info("git add *")
	err = w.AddWithOptions(&git.AddOptions{All: true})
	CheckIfError(err)

	Info("git status --porcelain")
	fmt.Println(w.Status())

	if s.Config.SkipCommit {
		Info("Skipping commit")
		return
	}
	Info("git commit -m \"New update from template\"")
	commit, err := w.Commit("New update from template", &git.CommitOptions{})

	obj, err := s.r.CommitObject(commit)
	fmt.Println(obj)
	CheckIfError(err)
}
