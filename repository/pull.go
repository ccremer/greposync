package repository

import (
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
)

func (s *Service) Pull() {
	if s.Config.SkipReset {
		s.p.WarnF("Skipped: pull")
		return
	}
	w, err := s.r.Worktree()
	s.p.CheckIfError(err)

	s.p.InfoF("git pull origin")
	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{})
	if err != git.NoErrAlreadyUpToDate {
		printer.CheckIfError(err)
	}
}
