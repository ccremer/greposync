package githosting

import (
	"fmt"

	"github.com/ccremer/greposync/domain"
)

type PullRequestStore struct {
	providers ProviderMap
}

func NewPullRequestStore(providers ProviderMap) *PullRequestStore {
	return &PullRequestStore{
		providers: providers,
	}
}

func (p *PullRequestStore) FindMatchingPullRequest(repository *domain.GitRepository) (*domain.PullRequest, error) {
	for _, remote := range p.providers {
		if remote.HasSupportFor(repository.URL) {
			pr, err := remote.FindPullRequest(repository.URL, repository.DefaultBranch, repository.CommitBranch)
			return pr, err
		}
	}
	return nil, fmt.Errorf("no remote providers supported: %s", repository.URL)
}

func (p *PullRequestStore) EnsurePullRequest(repository *domain.GitRepository) error {
	for _, remote := range p.providers {
		if remote.HasSupportFor(repository.URL) {
			return remote.EnsurePullRequest(repository.URL, repository.PullRequest)
		}
	}
	return fmt.Errorf("no remote providers supported: %s", repository.URL)
}
