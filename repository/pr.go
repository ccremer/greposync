package repository

import (
	"github.com/ccremer/git-repo-sync/repository/github"
	"github.com/go-git/go-git/v5"
)

func CreatePR(r *git.Repository) {
	github.CreatePR(r)
}
