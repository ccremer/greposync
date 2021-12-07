package valuestore

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStore_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*domain.ValueStore)(nil), new(MapStore))
}

func TestMapStore_loadAndMergeConfig(t *testing.T) {
	tests := map[string]struct {
		expectedConf  map[string]interface{}
		givenSyncFile string
	}{
		"GivenExistingFile_ThenLoadYaml": {
			givenSyncFile: "defaults.yml",
			expectedConf: map[string]interface{}{
				"object": map[string]interface{}{
					"key": "value",
				},
				"array": []interface{}{
					"string",
				},
			},
		},
		"GivenNonParseableFile_ThenIgnore": {
			givenSyncFile: "defaults.ini",
			expectedConf:  map[string]interface{}{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewMapStore(nil)
			u, err := url.Parse("https://github.com/ccremer/greposync")
			require.NoError(t, err)
			s.globalConfig = config{}
			SyncConfigFileName = tt.givenSyncFile
			repo := &domain.GitRepository{URL: domain.FromURL(u), RootDir: domain.NewFilePath("testdata")}
			result, err := s.loadAndMergeConfig(repo)
			require.NoError(t, err)
			assert.EqualValues(t, tt.expectedConf, result)
		})
	}
}

func TestMapStore_loadDataForTemplate(t *testing.T) {
	tests := map[string]struct {
		expectedConf          domain.Values
		givenSyncFile         string
		givenTemplateFileName string
	}{
		"GivenExistingSimpleFile_ThenLoadYaml": {
			givenSyncFile:         "sync.yml",
			givenTemplateFileName: "README.md",
			expectedConf: domain.Values{
				"title": "Hello World",
				"key":   "value",
			},
		},
		"GivenFileWithDirectoryDefaults_ThenLoadYaml": {
			givenSyncFile:         "advanced.yml",
			givenTemplateFileName: ".github/workflows/release.yaml",
			expectedConf: domain.Values{
				"title": "Hello World",
				"key":   "specific",
				"array": "overridden",
			},
		},
		"GivenConfigWithNestedObjects_WhenSinglePropertyOverridden_ThenMergeOnlyOverriddenProperty": {
			givenSyncFile:         "nested.yml",
			givenTemplateFileName: ".github/workflows/release.yaml",
			expectedConf: domain.Values{
				"features": map[string]interface{}{
					"testing": map[string]interface{}{
						"enabled":    false,
						"anotherKey": "value",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewMapStore(nil)
			cfg, err := s.loadYaml(filepath.Join("testdata", tt.givenSyncFile))
			require.NoError(t, err)
			result, err := s.loadValuesForTemplate(cfg, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result)
		})
	}
}
