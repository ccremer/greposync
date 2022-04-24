package valuestore

import (
	"path"
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var specialFlagsCases = map[string]struct {
	givenTemplateFileName string
	expectedBooleanFlag   bool
	expectedTargetPath    domain.Path
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

func TestFetchFilesToDelete(t *testing.T) {
	tests := map[string]struct {
		givenSyncFile     string
		givenTemplateList []*domain.Template
		expectedFiles     []domain.Path
	}{
		"GivenExistingFile_ThenExpectCorrectPropertyInheritance": {
			givenSyncFile: "specialflags.yml",
			expectedFiles: []domain.Path{
				"subdir/fileTrue",
				"topLevelFileTrue",
			},
		},
		"GivenExistingFile_WhenMergingWithGlobals_ThenExpectAllListedFiles": {
			givenSyncFile: "delete-globals.yml",
			expectedFiles: []domain.Path{
				"topLevelFileTrue",
				"subdir/fileTrue",
			},
		},
		"GivenTemplateList_WhenTemplateDeletedWithGlobals_ThenExpectTemplateInList": {
			givenSyncFile: "delete-globals.yml",
			givenTemplateList: []*domain.Template{
				{RelativePath: "a-file-not-in-sync.yml"},
			},
			expectedFiles: []domain.Path{
				"topLevelFileTrue",
				"subdir/fileTrue",
				"a-file-not-in-sync.yml",
			},
		},
		"GivenTemplateList_WhenTemplateExcludedViaHierarchy_ThenExpectNotInList": {
			givenSyncFile: "delete-globals.yml",
			givenTemplateList: []*domain.Template{
				{RelativePath: "subdir/a-file-not-in-sync.yml"},
			},
			expectedFiles: []domain.Path{
				"topLevelFileTrue",
				"subdir/fileTrue",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewKoanfStore(nil)
			k := koanf.New("")
			require.NoError(t, k.Load(file.Provider(path.Join("testdata", tt.givenSyncFile)), yaml.Parser()))
			result, err := s.loadFilesToDelete(k, tt.givenTemplateList)
			require.NoError(t, err)
			require.Len(t, result, len(tt.expectedFiles), "length of result list")
			for i := range tt.expectedFiles {
				assert.Contains(t, result, tt.expectedFiles[i], "expected file list")
			}
		})
	}
}

func TestLoadBooleanFlag(t *testing.T) {
	for _, flagName := range []string{"delete", "unmanaged"} {
		for name, tt := range specialFlagsCases {
			t.Run(name+"_With_"+flagName, func(t *testing.T) {
				s := NewKoanfStore(nil)
				k := koanf.New("")
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

func TestLoadTargetPath(t *testing.T) {
	for name, tt := range specialFlagsCases {
		t.Run(name, func(t *testing.T) {
			s := NewKoanfStore(nil)
			k := koanf.New("")
			require.NoError(t, k.Load(file.Provider(path.Join("testdata", "specialflags.yml")), yaml.Parser()))
			result, err := s.loadTargetPath(k, tt.givenTemplateFileName)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedTargetPath, result)
		})
	}
}
