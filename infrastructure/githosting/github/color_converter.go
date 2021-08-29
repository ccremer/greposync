package github

import (
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
)

type ColorConverter struct{}

func (ColorConverter) ConvertToEntity(color string) domain.Color {
	formatted := strings.ToUpper(fmt.Sprintf("#%s", color))
	return domain.Color(formatted)
}

func (ColorConverter) ConvertFromEntity(color domain.Color) *string {
	formatted := strings.ToLower(strings.TrimPrefix(color.String(), "#"))
	return &formatted
}
