package repository

import (
	"os"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/repository/github"
)

func (s *Service) CreatePR(config cfg.PullRequestConfig) {
	if !s.Config.CreatePR {
		s.p.WarnF("Skipped: Create PR")
		return
	}
	c := github.Config{
		Token:        os.Getenv("GITHUB_TOKEN"),
		Subject:      config.Subject,
		Repo:         s.Config.GetName(),
		RepoOwner:    s.Config.Namespace,
		CommitBranch: s.Config.CommitBranch,
		TargetBranch: config.TargetBranch,
		Body:         config.BodyTemplate,
	}
	github.CreatePR(c)
}
