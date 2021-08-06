package github

import (
	"github.com/ccremer/greposync/core"
)

type (
	// LabelConverter converts core.Label to cfg.RepositoryLabel and vice-versa
	LabelConverter struct{}
)

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertToEntity(labels []*LabelImpl) []core.Label {
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
func (LabelConverter) ConvertFromEntity(labels []core.Label) []*LabelImpl {
	if labels == nil || len(labels) == 0 {
		return []*LabelImpl{}
	}
	converted := make([]*LabelImpl, len(labels))
	for i := range labels {
		converted[i] = labels[i].(*LabelImpl)
	}
	return converted
}
