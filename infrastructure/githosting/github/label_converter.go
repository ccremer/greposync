package github

import (
	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v38/github"
)

// LabelConverter converts domain.Label to github.Label and vice-versa
type LabelConverter struct{}

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertToEntity(label *github.Label) domain.Label {
	if label == nil {
		return domain.Label{}
	}
	entity := domain.Label{
		Name:        label.GetName(),
		Description: label.GetDescription(),
	}
	color := ColorConverter{}.ConvertToEntity(label.GetColor())
	// there's no non-colored label on GitHub
	_ = entity.SetColor(color)
	return entity
}

// ConvertFromEntity converts the given object to another.
// Returns a nil if labels is empty or nil.
func (LabelConverter) ConvertFromEntity(label domain.Label) *github.Label {
	converted := &github.Label{
		Name:        &label.Name,
		Description: &label.Description,
	}
	color := ColorConverter{}.ConvertFromEntity(label.GetColor())
	converted.Color = &color
	return converted
}
