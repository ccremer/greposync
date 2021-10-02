package domain

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplate_CleanPath(t *testing.T) {
	tests := map[string]struct {
		givenPath    Path
		expectedPath Path
	}{
		"GivenFile_WhenNoExtension_ThenExpectSame": {
			givenPath:    NewPath("readme"),
			expectedPath: NewPath("readme"),
		},
		"GivenFile_WhenNoSpecialExtension_ThenExpectSame": {
			givenPath:    NewPath("readme.md"),
			expectedPath: NewPath("readme.md"),
		},
		"GivenFileInDir_WhenNoExtension_ThenExpectSame": {
			givenPath:    NewPath("dir", "readme.md"),
			expectedPath: NewPath("dir", "readme.md"),
		},
		"GivenFile_WhenExtension_ThenRemoveIt": {
			givenPath:    NewPath("readme.md.tpl"),
			expectedPath: NewPath("readme.md"),
		},
		"GivenFile_WhenExtensionTwice_ThenRemoveOne": {
			givenPath:    NewPath("readme.md.tpl.tpl"),
			expectedPath: NewPath("readme.md.tpl"),
		},
		"GivenFileInDir_WhenExtension_ThenRemoveFromFileName": {
			givenPath:    NewPath("dir", "readme.tpl.md"),
			expectedPath: NewPath("dir", "readme.md"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			template := NewTemplate(tt.givenPath, Permissions(0))
			result := template.CleanPath()
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

func TestTemplate_AsValues(t *testing.T) {
	file := "testdata/template-meta.tpl"
	info, err := os.Stat(file)
	require.NoError(t, err)
	subject := NewTemplate(NewPath(file), Permissions(info.Mode()))
	values := subject.AsValues()

	engine := DummyEngine{templatePath: file}
	result, err := engine.Execute(subject, values)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s\n%s\n", file, "0644"), result.String())
}
