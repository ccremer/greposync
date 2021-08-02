package repository

import (
	"net/url"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/githosting/github"
)

type (
	// Repository is the implementation for core.GitRepository.
	Repository struct {
		Config     *cfg.GitConfig
		coreConfig core.GitRepositoryConfig
		labels     []core.Label
	}
)

func NewGitRepository(cfg *cfg.GitConfig, labels []core.Label) *Repository {
	return &Repository{
		Config: cfg,
		coreConfig: core.GitRepositoryConfig{
			Provider: getProvider(cfg.Url),
			URL:      core.FromURL(cfg.Url),
		},
		labels: labels,
	}
}

// GetLabels implements core.GitRepository.
func (g *Repository) GetLabels() []core.Label {
	return g.labels
}

func getProvider(url *url.URL) core.GitHostingProvider {
	if url.Hostname() == "github.com" {
		return github.GitHubProviderKey
	}
	return ""
}
