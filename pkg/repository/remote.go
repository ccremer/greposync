package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type Remote interface {
	// FindLabels returns the same set of labels, but converted to core.Label and optionally a provider-specific reference attached.
	FindLabels(url *core.GitURL, labels []*cfg.RepositoryLabel) ([]core.Label, error)

	// FindPullRequest returns a remote-specific core.PullRequest or nil if it doesn't exist remotely.
	FindPullRequest(url *core.GitURL, config PullRequestProperties) (core.PullRequest, error)

	// NewPullRequest returns a new entity without creating it on the remote.
	NewPullRequest(url *core.GitURL, config PullRequestProperties) core.PullRequest
	EnsurePullRequest(url *core.GitURL, pr core.PullRequest) error
}
