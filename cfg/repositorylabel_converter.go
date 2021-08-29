package cfg

import (
	"errors"
	"fmt"

	"github.com/ccremer/greposync/domain"
)

type RepositoryLabelConverter struct{}

// ConvertToEntity converts the given object to another.
func (RepositoryLabelConverter) ConvertToEntity(label RepositoryLabel) (domain.Label, error) {
	if label.Name == "" {
		return domain.Label{}, errors.New("invalid label configuration: empty label name not allowed")
	}
	entity := domain.Label{
		Name:        label.Name,
		Description: label.Description,
	}
	err := entity.SetColor(domain.Color(label.Color))
	if err != nil {
		return entity, fmt.Errorf("invalid label configuration: invalid color for '%s': %w", label.Name, err)
	}
	return entity, nil
}
