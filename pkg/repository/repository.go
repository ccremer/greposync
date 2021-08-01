package repository

import (
	"net/url"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/githosting/github"
)

type (
	// GitRepository is the implementation for core.GitRepositoryFacade.
	GitRepository struct {
		Config     *cfg.GitConfig
		coreConfig core.GitRepositoryConfig
		labels     []core.GitRepositoryLabel
	}
)

func newGitRepositoryFacade(cfg *cfg.GitConfig, labels []core.GitRepositoryLabel) *GitRepository {
	return &GitRepository{
		Config: cfg,
		coreConfig: core.GitRepositoryConfig{
			Provider: getProvider(cfg.Url),
			URL:      core.FromURL(cfg.Url),
		},
		labels: labels,
	}
}

// GetLabels implements core.GitRepositoryFacade.
func (g *GitRepository) GetLabels() []core.GitRepositoryLabel {
	return g.labels
}

// GetConfig implements core.GitRepositoryFacade.
func (g *GitRepository) GetConfig() core.GitRepositoryConfig {
	return g.coreConfig
}

func getProvider(url *url.URL) core.GitHostingProvider {
	if url.Hostname() == "github.com" {
		return github.GitHubProviderKey
	}
	return ""
}
