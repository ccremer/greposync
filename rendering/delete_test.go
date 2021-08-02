package rendering

import (
	"path"
	"testing"
	"text/template"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_searchOrphanedFiles(t *testing.T) {
	tests := map[string]struct {
		givenValues       Values
		expectedFileNames []string
	}{
		"GivenExistingTemplateName_WhenSearching_ThenIgnore": {
			givenValues: Values{
				"Readme.md": map[string]interface{}{
					"delete": true,
				},
			},
			expectedFileNames: []string{},
		},
		"GivenNoMatchingTemplate_WhenSearching_ThenReturnOrphanedKeys": {
			givenValues: Values{
				".gitignore": map[string]interface{}{
					"delete": true,
				},
			},
			expectedFileNames: []string{".gitignore"},
		},
		"GivenDirectoriesInKeys_WhenSearching_ThenIgnoreDirectories": {
			givenValues: Values{
				"dir/": map[string]interface{}{
					"delete": true,
				},
				"dir/filename": map[string]interface{}{
					"delete": true,
				},
			},
			expectedFileNames: []string{"dir/filename"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				k:        koanf.New(""),
				instance: rendering.NewGoTemplateStore(nil),
			}
			r.instance.SetTemplateInstances(map[string]*template.Template{
				"Readme.md": nil,
			})
			require.NoError(t, r.k.Load(confmap.Provider(tt.givenValues, ""), nil))
			files := r.searchOrphanedFiles()
			assert.Equal(t, tt.expectedFileNames, files)
		})
	}
}

func (s *TemplateTestSuite) TestDeleteUnwantedFiles() {
	r := NewRenderer(&cfg.SyncConfig{
		Git: &cfg.GitConfig{
			Dir: s.SeedTargetDir,
		},
	}, koanf.New(""))
	s.Require().NoError(r.k.Load(confmap.Provider(Values{
		"readme.md": map[string]interface{}{
			"delete": true,
		},
	}, ""), nil))
	r.DeleteUnwantedFiles()()
	s.Assert().NoFileExists(path.Join(s.SeedTargetDir, "readme.md"))
}
