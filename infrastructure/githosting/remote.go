package githosting

import (
	"github.com/ccremer/greposync/domain"
)

type RemoteProvider string

type Remote interface {
	// FindLabels returns the same set of labels, but converted to core.Label and optionally a provider-specific reference attached.
	FindLabels(url *domain.GitURL) ([]*domain.Label, error)

	// FindPullRequest returns a remote-specific domain.PullRequest or nil if none matching the branches exist remotely.
	FindPullRequest(url *domain.GitURL, baseBranch, commitBranch string) (*domain.PullRequest, error)

	// EnsurePullRequest creates or updates the given domain.PullRequest.
	// The same rules as domain.PullRequestStore:EnsurePullRequest applies.
	EnsurePullRequest(url *domain.GitURL, pr *domain.PullRequest) error

	// HasSupportFor returns true if the remote implementation supports interacting with the remote API for the given repository URL.
	HasSupportFor(url *domain.GitURL) bool
}
