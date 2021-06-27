package repository

import (
	"os"

	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
)

func (s *Service) PrepareWorkspace() {
	gitDir := s.Config.GitDir
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		s.CloneGitRepository()
		s.CheckoutBranch()
		return
	}
	repo, err := git.PlainOpen(gitDir)
	s.p.CheckIfError(err)
	s.r = repo

	s.ResetRepository()
	s.CheckoutBranch()
	s.Pull()
}

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

func (s *Service) CloneGitRepository() {
	s.p.InfoF("git clone")
	gitDir := s.Config.GitDir
	repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
		URL:      s.Config.GitUrl,
		Progress: os.Stdout,
	})
	s.p.CheckIfError(err)
	s.r = repo
}
