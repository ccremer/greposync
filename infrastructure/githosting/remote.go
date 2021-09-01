package githosting

import (
	"github.com/ccremer/greposync/domain"
)

type ProviderMap map[RemoteProvider]Remote

type RemoteProvider string

type Remote interface {
	// FetchLabels returns the domain.LabelSet found for the given repository.
	// An empty set without error is returned if none found.
	FetchLabels(url *domain.GitURL) (domain.LabelSet, error)

	DeleteLabels(url *domain.GitURL, labels domain.LabelSet) error

	EnsureLabels(url *domain.GitURL, labels domain.LabelSet) error

	// FindPullRequest returns a remote-specific domain.PullRequest or nil if none matching the branches exist remotely.
	FindPullRequest(url *domain.GitURL, baseBranch, commitBranch string) (*domain.PullRequest, error)

	// EnsurePullRequest creates or updates the given domain.PullRequest.
	// The same rules as domain.PullRequestStore:EnsurePullRequest applies.
	EnsurePullRequest(url *domain.GitURL, pr *domain.PullRequest) error

	// HasSupportFor returns true if the remote implementation supports interacting with the remote API for the given repository URL.
	HasSupportFor(url *domain.GitURL) bool
}
