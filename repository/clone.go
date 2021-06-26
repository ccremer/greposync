package repository

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneGitRepository(url, dir string) error {
	gitDir := "repos/" + dir
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		_, err = git.PlainClone(gitDir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		return err
	}
	log.Println("git exists already")
	return nil
}
