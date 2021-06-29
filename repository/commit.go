package repository

import (
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
)

func (s *Service) MakeCommit() {
	if s.Config.SkipCommit {
		s.p.WarnF("Skipped: git commit")
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

	s.p.InfoF("git commit -m \"%s\"", s.Config.CommitMessage)
	commit, err := w.Commit(s.Config.CommitMessage, &git.CommitOptions{})

	obj, err := s.r.CommitObject(commit)
	s.p.CheckIfError(err)
	s.p.UseColor(printer.White).LogF(obj.String())
}
