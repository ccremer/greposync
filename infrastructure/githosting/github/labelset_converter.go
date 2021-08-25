package github

import (
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v38/github"
)

type (
	// LabelConverter converts core.Label to cfg.RepositoryLabel and vice-versa
	LabelConverter struct{}

	ColorConverter struct{}
)

// ConvertToEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertToEntity(labels []*github.Label) domain.LabelSet {
	if labels == nil || len(labels) == 0 {
		return domain.LabelSet{}
	}
	converted := make([]domain.Label, len(labels))
	for i := range labels {
		original := labels[i]
		entity := domain.Label{
			Name:        *original.Name,
			Description: *original.Description,
		}
		color := ColorConverter{}.ConvertToEntity(*original.Color)
		// there's no non-colored label on GitHub
		_ = entity.SetColor(color)
		converted[i] = entity
	}
	return converted
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (LabelConverter) ConvertFromEntity(labels domain.LabelSet) []*github.Label {
	if labels == nil || len(labels) == 0 {
		return []*github.Label{}
	}
	converted := make([]*github.Label, len(labels))
	for i := range labels {
		original := labels[i]
		ghLabel := &github.Label{
			Name:        &original.Name,
			Color:       ColorConverter{}.ConvertFromEntity(original.GetColor()),
			Description: &original.Description,
		}
		converted[i] = ghLabel
	}
	return converted
}


func (ColorConverter) ConvertToEntity(color string) domain.Color {
	formatted := strings.ToUpper(fmt.Sprintf("#%s", color))
	return domain.Color(formatted)
}

func (ColorConverter) ConvertFromEntity(color domain.Color) *string {
	formatted := strings.ToLower(strings.TrimPrefix(color.String(), "#"))
	return &formatted
}
