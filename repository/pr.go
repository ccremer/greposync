package repository

import (
	"os"

	"github.com/ccremer/git-repo-sync/repository/github"
)

func (s *Service) CreatePR() {
	if !s.Config.CreatePR {
		Info("Skip creating PR")
		return
	}
	c := github.Config{
		Token:        os.Getenv("GITHUB_TOKEN"),
		Subject:      "Update from gsync",
		Repo:         "git-repo-sync",
		RepoOwner:    "ccremer",
		CommitBranch: "my-branch",
		TargetBranch: "master",
		Body:         "long text for PR desc",
	}
	github.CreatePR(c)
}
