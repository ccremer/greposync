package repository

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneGitRepository(url, dir string) (*git.Repository, error) {
	gitDir := "repos/" + dir
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		repo, err := git.PlainClone(gitDir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		return repo, err
	}
	return git.PlainOpen(gitDir)
}
