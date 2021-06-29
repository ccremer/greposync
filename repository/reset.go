package repository

import (
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
)

func (s *Service) ResetRepository() {
	if s.Config.SkipReset {
		s.p.WarnF("Skipped: git reset")
		return
	}
	s.p.InfoF("git fetch")
	err := s.r.Fetch(&git.FetchOptions{})
	if err != git.NoErrAlreadyUpToDate {
		printer.CheckIfError(err)
	}

	w, err := s.r.Worktree()
	s.p.CheckIfError(err)

	s.p.InfoF("git reset --hard")
	err = w.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	s.p.CheckIfError(err)
}
