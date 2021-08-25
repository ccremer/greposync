package githosting

import (
	"fmt"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
)

type PullRequestStore struct {
	remote    Remote
	providers ProviderMap
	config    *cfg.Configuration
}

type ProviderMap map[RemoteProvider]Remote

func NewPullRequestStore(config *cfg.Configuration, providers ProviderMap) *PullRequestStore {
	return &PullRequestStore{
		providers: providers,
		config:    config,
	}
}

func (p *PullRequestStore) FindMatchingPullRequest(repository *domain.GitRepository) (*domain.PullRequest, error) {
	for _, remote := range p.providers {
		if remote.HasSupportFor(repository.URL) {
			pr, err := p.remote.FindPullRequest(repository.URL, repository.DefaultBranch, repository.CommitBranch)
			return pr, err
		}
	}
	return nil, fmt.Errorf("no remote provider supported: %s", repository.URL)
}

func (p *PullRequestStore) EnsurePullRequest(repository *domain.GitRepository) error {
	for _, remote := range p.providers {
		if remote.HasSupportFor(repository.URL) {
			return p.remote.EnsurePullRequest(repository.URL, repository.PullRequest)
		}
	}
	return fmt.Errorf("no remote provider supported: %s", repository.URL)
}
