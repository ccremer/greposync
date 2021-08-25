package repository

import (
	"fmt"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/ccremer/greposync/repository"
)

// RepositoryStore is the implementation of core.GitRepositoryStore.
type RepositoryStore struct {
	providers ProviderMap
	config    *cfg.Configuration
	cache     map[*core.GitURL]core.GitRepository
}

type ProviderMap map[core.GitHostingProvider]githosting.Remote

// NewRepositoryStore creates a new instance.
func NewRepositoryStore(config *cfg.Configuration, providers ProviderMap) *RepositoryStore {
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
		if provider, exists := r.providers[service.Config.Provider]; exists {
			labels, err = provider.FindLabels(core.FromURL(service.Config.Url), convertLabels(r.config.RepositoryLabels))
			if err != nil {
				return []core.GitRepository{}, err
			}
		}
		repos[i] = NewGitRepository(service.Config, r.config.PullRequest, labels)
	}
	return repos, err
}

// FetchGitRepository implements core.GitRepositoryStore.
func (r *RepositoryStore) FetchGitRepository(url *core.GitURL) (core.GitRepository, error) {
	if r.cache == nil {
		panic("implement me")
	}
	if repo, exists := r.cache[url]; exists {
		return repo, nil
	}
	return nil, fmt.Errorf("repository with url '%s' not found", url.Redacted())
}

func convertLabels(labels map[string]cfg.RepositoryLabel) []*cfg.RepositoryLabel {
	values := make([]*cfg.RepositoryLabel, 0)
	for _, label := range labels {
		values = append(values, &label)
	}
	return values
}
