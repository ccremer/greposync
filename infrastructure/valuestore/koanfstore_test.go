package valuestore

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKoanfStore_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*domain.ValueStore)(nil), new(KoanfStore))
}

func TestKoanfStore_loadAndMergeConfig(t *testing.T) {
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
			s := NewKoanfStore(nil)
			s.syncConfigFileName = tt.givenSyncFile
			u, err := url.Parse("https://github.com/ccremer/greposync")
			require.NoError(t, err)
			s.globalKoanf = koanf.New("")
			repo := &domain.GitRepository{URL: domain.FromURL(u), RootDir: domain.NewFilePath("testdata")}
			result, err := s.loadAndMergeConfig(repo)
			require.NoError(t, err)
			assert.EqualValues(t, tt.expectedConf, result.Raw())
		})
	}
}

func TestKoanfStore_LoadDataForTemplate(t *testing.T) {
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
			s := NewKoanfStore(nil)
			k := koanf.New("")
			err := s.loadYaml(k, filepath.Join("testdata", tt.givenSyncFile))
			require.NoError(t, err)
			result, err := s.loadValuesForTemplate(k, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result)
		})
	}
}

func TestKoanfStore_FetchValuesForTemplate(t *testing.T) {
	newRepo := func(t *testing.T, path string) *domain.GitRepository {
		u, err := url.Parse(path)
		require.NoError(t, err)
		return domain.NewGitRepository(domain.FromURL(u), domain.NewFilePath(path))
	}
	t.Run("GivenMultipleRepos_ThenSyncFilesAreIsolated", func(t *testing.T) {
		s := NewKoanfStore(nil)
		tpl := domain.NewTemplate(domain.NewFilePath("README.md"), 0777)
		s.globalKoanf = koanf.New(".")
		err := s.globalKoanf.Load(confmap.Provider(config{
			"README.md": map[string]interface{}{
				"global": "value",
			},
		}, ""), nil)
		require.NoError(t, err)
		repo1 := newRepo(t, "testdata/repo1")
		repo2 := newRepo(t, "testdata/repo2")
		vals1, err := s.FetchValuesForTemplate(tpl, repo1)
		require.NoError(t, err)
		vals2, err := s.FetchValuesForTemplate(tpl, repo2)
		require.NoError(t, err)

		assert.Equal(t, "parameter", vals1["extra"])
		assert.Nil(t, vals2["extra"])
	})
}
