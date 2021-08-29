package valuestore

import (
	"path/filepath"
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKoanfValueStore_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*domain.ValueStore)(nil), new(KoanfValueStore))
}

func TestKoanfValueStore_loadAndMergeConfig(t *testing.T) {
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
			s := NewValueStore()
			s.globalKoanf = koanf.New("")
			result, err := s.loadAndMergeConfig(filepath.Join("testdata", tt.givenSyncFile))
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result.Raw())
		})
	}
}

func TestKoanfValueStore_loadDataForTemplate(t *testing.T) {
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
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewValueStore()
			k := koanf.New(".")
			require.NoError(t, k.Load(file.Provider(filepath.Join("testdata", tt.givenSyncFile)), yaml.Parser()))
			result, err := s.loadDataForTemplate(k, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result)
		})
	}
}
