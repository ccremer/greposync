package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/githosting/github"
	"github.com/ccremer/greposync/repository"
)

type (
	// GitRepositoryProvider is the implementation of core.ManagedRepoProvider.
	GitRepositoryProvider struct {
		providers map[core.GitHostingProvider]core.GitHostingFacade
		config    *cfg.Configuration
	}
)

// NewGitRepositoryProvider creates a new instance.
func NewGitRepositoryProvider(config *cfg.Configuration) *GitRepositoryProvider {
	return &GitRepositoryProvider{
		providers: getProviders(),
		config:    config,
	}
}

// GetSupportedGitHostingProviders implements core.ManagedRepoProvider.
func (r *GitRepositoryProvider) GetSupportedGitHostingProviders() map[core.GitHostingProvider]core.GitHostingFacade {
	return r.providers
}

// LoadManagedRepositories implements core.ManagedRepoProvider.
func (r *GitRepositoryProvider) LoadManagedRepositories() ([]core.GitRepositoryFacade, error) {
	services, err := repository.NewServicesFromFile(r.config)
	if err != nil || len(services) == 0 {
		return []core.GitRepositoryFacade{}, err
	}
	labels := convertLabels(r.config.RepositoryLabels)
	facades := make([]core.GitRepositoryFacade, len(services))
	for i, service := range services {
		facade := newGitRepositoryFacade(service.Config, labels)
		facades[i] = facade
	}
	return facades, err
}

func getProviders() map[core.GitHostingProvider]core.GitHostingFacade {
	return map[core.GitHostingProvider]core.GitHostingFacade{
		github.GitHubProviderKey: github.NewFacade(),
	}
}

func convertLabels(labels map[string]cfg.RepositoryLabel) []core.GitRepositoryLabel {
	values := make([]*cfg.RepositoryLabel, 0)
	for _, label := range labels {
		values = append(values, &label)
	}
	return github.LabelConverter{}.ConvertToEntity(values)
}
