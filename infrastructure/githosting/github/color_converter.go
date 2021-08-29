package github

import (
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
)

type ColorConverter struct{}

func (ColorConverter) ConvertToEntity(color string) domain.Color {
	formatted := strings.ToUpper(fmt.Sprintf("#%s", color))
	converted := domain.Color(formatted)
	err := converted.CheckValue()
	if err != nil {
		return ""
	}
	return converted
}

func (ColorConverter) ConvertFromEntity(color domain.Color) string {
	err := color.CheckValue()
	if err != nil {
		return ""
	}
	formatted := strings.ToLower(strings.TrimPrefix(color.String(), "#"))
	return formatted
}
