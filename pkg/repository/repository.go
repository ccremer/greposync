package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type (
	// Repository is the implementation for core.GitRepository.
	Repository struct {
		GitConfig  *cfg.GitConfig
		PrConfig   *cfg.PullRequestConfig
		coreConfig core.GitRepositoryConfig
		labels     []core.Label
		remote     Remote
		pr         core.PullRequest
	}
)

func NewGitRepository(cfg *cfg.GitConfig, prConfig *cfg.PullRequestConfig, labels []core.Label) *Repository {
	return &Repository{
		GitConfig: cfg,
		coreConfig: core.GitRepositoryConfig{
			URL:     core.FromURL(cfg.Url),
			RootDir: cfg.Dir,
		},
		PrConfig: prConfig,
		labels:   labels,
	}
}

// GetLabels implements core.GitRepository.
func (g *Repository) GetLabels() []core.Label {
	return g.labels
}
