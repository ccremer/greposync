package rendering

import (
	"testing"
	"text/template"

	"github.com/ccremer/greposync/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplate_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*core.Template)(nil), new(GoTemplate))
}

func TestGoTemplate_Render(t *testing.T) {
	tpl, err := template.New("").Funcs(GoTemplateFuncMap()).Parse("{{ .data | toJson | upper }}")
	require.NoError(t, err)
	gotemplate := GoTemplate{template: tpl}
	result, err := gotemplate.Render(core.Values{
		"data": core.Values{
			"json": true,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "{\"JSON\":TRUE}", result)
}

func TestGoTemplate_GetRelativePath(t *testing.T) {
	tests := map[string]struct {
		givenPath    string
		expectedPath string
	}{
		"GivenFileWithoutDir_WhenSanitizing_ThenReturnSamePath": {
			givenPath:    "fileName",
			expectedPath: "fileName",
		},
		"GivenFileInDir_WhenSanitizing_ThenReturnSamePath": {
			givenPath:    "dir/fileName",
			expectedPath: "dir/fileName",
		},
		"GivenFileWithTplExtension_WhenSanitizing_ThenReturnStripped": {
			givenPath:    "dir/fileName.tpl",
			expectedPath: "dir/fileName",
		},
		"GivenFileWithTplExtensionTwice_WhenSanitizing_ThenReturnStrippedOnce": {
			givenPath:    "fileName.tpl.tpl",
			expectedPath: "fileName.tpl",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &GoTemplate{
				RelativePath: tt.givenPath,
			}
			result := g.GetRelativePath()
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}
