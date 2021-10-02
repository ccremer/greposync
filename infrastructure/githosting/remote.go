package githosting

import (
	"github.com/ccremer/greposync/domain"
)

type ProviderMap map[RemoteProvider]Remote

type RemoteProvider string

type Remote interface {
	// FetchLabels returns the domain.LabelSet found for the given repository.
	// An empty set without error is returned if none found.
	FetchLabels(repository *domain.GitRepository) (domain.LabelSet, error)

	DeleteLabels(repository *domain.GitRepository, labels domain.LabelSet) error

	EnsureLabels(repository *domain.GitRepository, labels domain.LabelSet) error

	// FindPullRequest returns a remote-specific domain.PullRequest or nil if none matching the branches exist remotely.
	FindPullRequest(repository *domain.GitRepository) (*domain.PullRequest, error)

	// EnsurePullRequest creates or updates the given domain.PullRequest.
	// The same rules as domain.PullRequestStore:EnsurePullRequest applies.
	EnsurePullRequest(repository *domain.GitRepository, pr *domain.PullRequest) error

	// HasSupportFor returns true if the remote implementation supports interacting with the remote API for the given repository URL.
	HasSupportFor(url *domain.GitURL) bool
}
