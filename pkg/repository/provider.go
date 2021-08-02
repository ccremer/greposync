package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/githosting/github"
	"github.com/ccremer/greposync/repository"
)

type (
	// RepositoryStore is the implementation of core.GitRepositoryStore.
	RepositoryStore struct {
		providers map[core.GitHostingProvider]core.GitHostingFacade
		config    *cfg.Configuration
	}
)

// NewRepositoryStore creates a new instance.
func NewRepositoryStore(config *cfg.Configuration) *RepositoryStore {
	return &RepositoryStore{
		providers: getProviders(),
		config:    config,
	}
}

// GetSupportedGitHostingProviders implements core.GitRepositoryStore.
func (r *RepositoryStore) GetSupportedGitHostingProviders() map[core.GitHostingProvider]core.GitHostingFacade {
	return r.providers
}

// FetchGitRepositories implements core.GitRepositoryStore.
func (r *RepositoryStore) FetchGitRepositories() ([]core.GitRepository, error) {
	services, err := repository.NewServicesFromFile(r.config)
	if err != nil || len(services) == 0 {
		return []core.GitRepository{}, err
	}
	labels := convertLabels(r.config.RepositoryLabels)
	facades := make([]core.GitRepository, len(services))
	for i, service := range services {
		facade := NewGitRepository(service.Config, labels)
		facades[i] = facade
	}
	return facades, err
}

func getProviders() map[core.GitHostingProvider]core.GitHostingFacade {
	return map[core.GitHostingProvider]core.GitHostingFacade{
		github.GitHubProviderKey: github.NewFacade(),
	}
}

func convertLabels(labels map[string]cfg.RepositoryLabel) []core.Label {
	values := make([]*cfg.RepositoryLabel, 0)
	for _, label := range labels {
		values = append(values, &label)
	}
	return github.LabelConverter{}.ConvertToEntity(values)
}
