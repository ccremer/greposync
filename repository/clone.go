package repository

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func (s *Service) PrepareWorkspace() {
	if _, err := os.Stat(s.Config.Dir); os.IsNotExist(err) {
		s.CloneGitRepository()
		s.CheckoutBranch()
		return
	}
	repo, err := git.PlainOpen(s.Config.Dir)
	s.p.CheckIfError(err)
	s.r = repo

	s.ResetRepository()
	s.CheckoutBranch()
	s.Pull()
}

func (s *Service) CloneGitRepository() {
	s.p.InfoF("git clone")
	gitDir := s.Config.Dir
	repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
		URL:      s.Config.Url,
		Progress: os.Stdout,
	})
	s.p.CheckIfError(err)
	s.r = repo
}
