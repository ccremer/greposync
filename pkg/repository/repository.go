package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type (
	// Repository is the implementation for core.GitRepository.
	Repository struct {
		Config     *cfg.GitConfig
		coreConfig core.GitRepositoryConfig
		labels     []core.Label
		remote     Remote
	}
)

func NewGitRepository(cfg *cfg.GitConfig, labels []core.Label) *Repository {
	return &Repository{
		Config: cfg,
		coreConfig: core.GitRepositoryConfig{
			URL: core.FromURL(cfg.Url),
		},
		labels: labels,
	}
}

// GetLabels implements core.GitRepository.
func (g *Repository) GetLabels() []core.Label {
	return g.labels
}
