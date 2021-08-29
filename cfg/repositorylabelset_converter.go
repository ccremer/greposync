package cfg

import (
	"github.com/ccremer/greposync/domain"
)

type RepositoryLabelSetConverter struct{}

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (RepositoryLabelSetConverter) ConvertToEntity(labels []RepositoryLabel) (domain.LabelSet, error) {
	if labels == nil || len(labels) == 0 {
		return domain.LabelSet{}, nil
	}
	converted := make(domain.LabelSet, len(labels))
	for i := range labels {
		label := labels[i]
		entity, err := RepositoryLabelConverter{}.ConvertToEntity(label)
		if err != nil {
			return converted, err
		}
		converted[i] = entity
	}
	return converted, nil
}
