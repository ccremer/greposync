package github

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type (
	// LabelConverter converts core.GitRepositoryLabel to cfg.RepositoryLabel and vice-versa
	LabelConverter struct{}
)

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertToEntity(labels []*cfg.RepositoryLabel) []core.GitRepositoryLabel {
	if labels == nil || len(labels) == 0 {
		return []core.GitRepositoryLabel{}
	}
	converted := make([]core.GitRepositoryLabel, len(labels))
	for i := range labels {
		converted[i] = labels[i]
	}
	return converted
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertFromEntity(labels []core.GitRepositoryLabel) []*cfg.RepositoryLabel {
	if labels == nil || len(labels) == 0 {
		return []*cfg.RepositoryLabel{}
	}
	converted := make([]*cfg.RepositoryLabel, len(labels))
	for i := range labels {
		converted[i] = labels[i].(*cfg.RepositoryLabel)
	}
	return converted
}
