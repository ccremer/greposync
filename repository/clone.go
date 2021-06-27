package repository

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func (s *Service) PrepareWorkspace() {
	gitDir := s.Config.GitDir
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		s.CloneGitRepository()
		s.SwitchBranch()
		return
	}
	repo, err := git.PlainOpen(gitDir)
	CheckIfError(err)

	//ResetRepository(repo)
	//SwitchBranch(repo)
	//Pull(repo)

	s.r = repo
}

func (s *Service) ResetRepository() {
	Info("git fetch")
	err := s.r.Fetch(&git.FetchOptions{})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	w, err := s.r.Worktree()
	CheckIfError(err)

	Info("git reset --hard")
	err = w.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	CheckIfError(err)
}

func (s *Service) CloneGitRepository() {
	gitDir := s.Config.GitDir
	repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
		URL:      s.Config.GitUrl,
		Progress: os.Stdout,
	})
	CheckIfError(err)
	s.r = repo
}
