package repository

import (
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
)

func (s *Service) MakeCommit() {
	if s.Config.SkipCommit {
		return
	}
	w, err := s.r.Worktree()
	s.p.CheckIfError(err)

	s.p.InfoF("git add *")
	err = w.AddWithOptions(&git.AddOptions{All: true})
	s.p.CheckIfError(err)

	s.p.InfoF("git status --porcelain")
	status, err := w.Status()
	s.p.CheckIfError(err)
	s.p.LogF(status.String())

	if s.Config.SkipCommit {
		s.p.WarnF("Skipped: commit")
		return
	}
	s.p.InfoF("git commit -m \"New update from template\"")
	commit, err := w.Commit("New update from template", &git.CommitOptions{})

	obj, err := s.r.CommitObject(commit)
	s.p.CheckIfError(err)
	s.p.UseColor(printer.White).LogF(obj.String())
}
