package repository

import (
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

func PrepareWorkspace(url, dir string) *git.Repository {
	gitDir := path.Join("repos", dir)
	gitDir = path.Clean(gitDir)
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		repo := CloneGitRepository(url, dir)
		SwitchBranch(repo)
		return repo
	}
	repo, err := git.PlainOpen(gitDir)
	CheckIfError(err)

	ResetRepository(repo)
	SwitchBranch(repo)
	Pull(repo)

	return repo
}

func ResetRepository(repo *git.Repository) {
	Info("git fetch")
	err := repo.Fetch(&git.FetchOptions{})
	if err != git.NoErrAlreadyUpToDate {
		CheckIfError(err)
	}

	w, err := repo.Worktree()
	CheckIfError(err)

	Info("git reset --hard")
	err = w.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	CheckIfError(err)
}

func CloneGitRepository(url, dir string) *git.Repository {
	gitDir := path.Join("repos", dir)
	gitDir = path.Clean(gitDir)
	repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	CheckIfError(err)
	return repo
}
