package valuestore

import (
	"path"
	"testing"

	"github.com/ccremer/greposync/core"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKoanfValueStore_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*core.ValueStore)(nil), new(KoanfValueStore))
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
			s := NewValueStore(koanf.New("."))
			result, err := s.loadAndMergeConfig(path.Join("testdata", tt.givenSyncFile))
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result.Raw())
		})
	}
}

func TestKoanfValueStore_loadDataForTemplate(t *testing.T) {
	tests := map[string]struct {
		expectedConf          core.Values
		givenSyncFile         string
		givenTemplateFileName string
	}{
		"GivenExistingSimpleFile_ThenLoadYaml": {
			givenSyncFile:         "sync.yml",
			givenTemplateFileName: "README.md",
			expectedConf: core.Values{
				"title": "Hello World",
				"key":   "value",
			},
		},
		"GivenFileWithDirectoryDefaults_ThenLoadYaml": {
			givenSyncFile:         "advanced.yml",
			givenTemplateFileName: ".github/workflows/release.yaml",
			expectedConf: core.Values{
				"title": "Hello World",
				"key":   "specific",
				"array": "overridden",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewValueStore(koanf.New("."))
			k := koanf.New(".")
			require.NoError(t, k.Load(file.Provider(path.Join("testdata", tt.givenSyncFile)), yaml.Parser()))
			result, err := s.loadDataForTemplate(k, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedConf, result)
		})
	}
}

var specialFlagsCases = map[string]struct {
	givenTemplateFileName string
	expectedBooleanFlag   bool
	expectedTargetPath    string
	expectedErrString     string
}{
	"GivenTopLevelFile_WhenTrue_ThenExpectTrue": {
		givenTemplateFileName: "topLevelFileTrue",
		expectedBooleanFlag:   true,
		expectedTargetPath:    "topLevelFile",
	},
	"GivenTopLevelFile_WhenFalse_ThenExpectFalse": {
		givenTemplateFileName: "topLevelFileFalse",
		expectedBooleanFlag:   false,
		expectedTargetPath:    "movedToDir/topLevelFileFalse",
	},
	"GivenFileInSubdir_WhenTrueInheritedByDir_ThenExpectTrue": {
		givenTemplateFileName: "subdir/fileTrue",
		expectedBooleanFlag:   true,
		expectedTargetPath:    "movedDir/fileTrue",
	},
	"GivenFileInSubdir_WhenFalseSetExplicitly_ThenExpectFalse": {
		givenTemplateFileName: "subdir/fileFalse",
		expectedBooleanFlag:   false,
		expectedTargetPath:    "anotherDir/file.renamed",
	},
	"GivenTopLevelFile_WhenInvalidDataType_ThenExpectFalse": {
		givenTemplateFileName: "invalidFile",
		expectedBooleanFlag:   false,
	},
	"GivenUndefinedFile_ThenExpectError": {
		givenTemplateFileName: "undefined",
		expectedBooleanFlag:   false,
		expectedErrString:     "key not found",
	},
}

func TestKoanfValueStore_FetchFilesToDelete(t *testing.T) {
	tests := map[string]struct {
		givenSyncFile string
		expectedFiles []string
	}{
		"GivenExistingFile_ThenExpectCorrectPropertyInheritance": {
			givenSyncFile: "specialflags.yml",
			expectedFiles: []string{
				"subdir/fileTrue",
				"topLevelFileTrue",
			},
		},
		"GivenExistingFile_WhenMergingWithGlobals_ThenExpectAllListedFiles": {
			givenSyncFile: "delete-globals.yml",
			expectedFiles: []string{
				"topLevelFileTrue",
				"subdir/fileTrue",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewValueStore(koanf.New("."))
			k := koanf.New(".")
			require.NoError(t, k.Load(file.Provider(path.Join("testdata", tt.givenSyncFile)), yaml.Parser()))
			result, err := s.loadFilesToDelete(k)
			require.NoError(t, err)
			require.Len(t, result, len(tt.expectedFiles))
			for i := range tt.expectedFiles {
				assert.Contains(t, result, tt.expectedFiles[i])
			}
		})
	}
}

func TestKoanfValueStore_loadBooleanFlag(t *testing.T) {
	for _, flagName := range []string{"delete", "unmanaged"} {
		for name, tt := range specialFlagsCases {
			t.Run(name+"_With_"+flagName, func(t *testing.T) {
				s := NewValueStore(koanf.New("."))
				k := koanf.New(".")
				require.NoError(t, k.Load(file.Provider(path.Join("testdata", "specialflags.yml")), yaml.Parser()))
				result, err := s.loadBooleanFlag(k, tt.givenTemplateFileName, flagName)
				if tt.expectedErrString != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErrString)
				} else {
					require.NoError(t, err)
				}
				assert.Equal(t, tt.expectedBooleanFlag, result)
			})
		}
	}
}

func TestKoanfValueStore_FetchTargetPath(t *testing.T) {
	for name, tt := range specialFlagsCases {
		t.Run(name, func(t *testing.T) {
			s := NewValueStore(koanf.New("."))
			k := koanf.New(".")
			require.NoError(t, k.Load(file.Provider(path.Join("testdata", "specialflags.yml")), yaml.Parser()))
			result, err := s.loadTargetPath(k, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedTargetPath, result)
		})
	}
}
