package github

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

type (
	// LabelConverter converts core.Label to cfg.RepositoryLabel and vice-versa
	LabelConverter struct{}
)

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertToEntity(labels []*cfg.RepositoryLabel) []core.Label {
	if labels == nil || len(labels) == 0 {
		return []core.Label{}
	}
	converted := make([]core.Label, len(labels))
	for i := range labels {
		converted[i] = labels[i]
	}
	return converted
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertFromEntity(labels []core.Label) []*cfg.RepositoryLabel {
	if labels == nil || len(labels) == 0 {
		return []*cfg.RepositoryLabel{}
	}
	converted := make([]*cfg.RepositoryLabel, len(labels))
	for i := range labels {
		converted[i] = labels[i].(*cfg.RepositoryLabel)
	}
	return converted
}
