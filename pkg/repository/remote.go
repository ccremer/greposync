package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type Remote interface {
	// Initialize prepares the remote API for usage.
	// It is only called once at program startup.
	Initialize() error
	// FindLabels returns the same set of labels, but converted to core.Label and optionally a provider-specific reference attached.
	FindLabels(url *core.GitURL, labels []*cfg.RepositoryLabel) ([]core.Label, error)
}
