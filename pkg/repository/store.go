package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/repository"
)

type (
	// RepositoryStore is the implementation of core.GitRepositoryStore.
	RepositoryStore struct {
		providers map[core.GitHostingProvider]Remote
		config    *cfg.Configuration
	}
)

// NewRepositoryStore creates a new instance.
func NewRepositoryStore(config *cfg.Configuration, providers map[core.GitHostingProvider]Remote) *RepositoryStore {
	return &RepositoryStore{
		providers: providers,
		config:    config,
	}
}

// FetchGitRepositories implements core.GitRepositoryStore.
func (r *RepositoryStore) FetchGitRepositories() ([]core.GitRepository, error) {
	services, err := repository.NewServicesFromFile(r.config)
	if err != nil || len(services) == 0 {
		return []core.GitRepository{}, err
	}
	repos := make([]core.GitRepository, len(services))
	for i, service := range services {

		var labels = make([]core.Label, 0)
		if provider, exists := r.providers[core.GitHostingProvider(service.Config.Provider)]; exists {
			labels, err = provider.FindLabels(core.FromURL(service.Config.Url), convertLabels(r.config.RepositoryLabels))
			if err != nil {
				return []core.GitRepository{}, err
			}
		}
		repos[i] = NewGitRepository(service.Config, labels)
	}
	return repos, err
}

func convertLabels(labels map[string]cfg.RepositoryLabel) []*cfg.RepositoryLabel {
	values := make([]*cfg.RepositoryLabel, 0)
	for _, label := range labels {
		values = append(values, &label)
	}
	return values
}
