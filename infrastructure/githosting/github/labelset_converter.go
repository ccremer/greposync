package github

import (
	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v39/github"
)

// LabelSetConverter converts domain.Label to github.Label and vice-versa
type LabelSetConverter struct{}

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelSetConverter) ConvertToEntity(labels []*github.Label) domain.LabelSet {
	if labels == nil || len(labels) == 0 {
		return domain.LabelSet{}
	}
	converted := make(domain.LabelSet, len(labels))
	for i := range labels {
		label := labels[i]
		entity := LabelConverter{}.ConvertToEntity(label)
		converted[i] = entity
	}
	return converted
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelSetConverter) ConvertFromEntity(labels domain.LabelSet) []*github.Label {
	if labels == nil || len(labels) == 0 {
		return []*github.Label{}
	}
	converted := make([]*github.Label, len(labels))
	for i := range labels {
		label := labels[i]
		ghLabel := LabelConverter{}.ConvertFromEntity(label)
		converted[i] = ghLabel
	}
	return converted
}
