package repository

import (
	"os"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/repository/github"
	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) CreateOrUpdatePR(config *cfg.PullRequestConfig) pipeline.ActionFunc {
	return func() pipeline.Result {
		if config.TargetBranch == "" {
			config.TargetBranch = s.Config.DefaultBranch
		}
		c := &github.Config{
			Token:        os.Getenv("GITHUB_TOKEN"),
			Subject:      config.Subject,
			Repo:         s.Config.GetName(),
			RepoOwner:    s.Config.Namespace,
			CommitBranch: s.Config.CommitBranch,
			TargetBranch: config.TargetBranch,
			Body:         config.BodyTemplate,
		}
		gh := github.NewProvider(c)
		return pipeline.Result{Err: gh.CreateOrUpdatePR()}
	}
}
