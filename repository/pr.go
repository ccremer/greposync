package repository

import (
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/repository/github"
)

// CreateOrUpdatePr creates a PR if it doesn't exist or updates if the remote branch exists already.
func (s *Service) CreateOrUpdatePr(config *cfg.PullRequestConfig) pipeline.ActionFunc {
	return func() pipeline.Result {
		if config.TargetBranch == "" {
			config.TargetBranch = s.Config.DefaultBranch
		}
		c := &github.Config{
			Token:        os.Getenv("GITHUB_TOKEN"),
			Subject:      config.Subject,
			Repo:         s.Config.Name,
			RepoOwner:    s.Config.Namespace,
			CommitBranch: s.Config.CommitBranch,
			TargetBranch: config.TargetBranch,
			Body:         config.BodyTemplate,
			Labels:       config.Labels,
		}
		gh := github.NewProvider(c)
		return pipeline.Result{Err: gh.CreateOrUpdatePr()}
	}
}
