package rendering

import (
	"testing"

	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_loadVariables(t *testing.T) {
	tests := map[string]struct {
		givenSyncFile  string
		givenGlobals   Values
		expectErr      bool
		expectedValues Values
	}{
		"GivenSyncFileNotExisting_WhenLoading_ThenUseGlobalVars": {
			givenSyncFile:  "ignored",
			givenGlobals:   Values{"Key": "Value"},
			expectedValues: Values{"Key": "Value"},
		},
		"GivenSyncExists_WhenLoading_ThenOverrideGlobalVars": {
			givenSyncFile:  "testdata/.sync.yml",
			givenGlobals:   Values{":globals": Values{"key": "value"}},
			expectedValues: Values{":globals": map[string]interface{}{"key": "overridden", "new": "variable"}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				p:              printer.DefaultPrinter,
				k:              koanf.New("."),
				globalDefaults: koanf.New("."),
			}
			require.NoError(t, r.globalDefaults.Load(confmap.Provider(tt.givenGlobals, ""), nil))
			err := r.loadVariables(tt.givenSyncFile)
			if tt.expectErr {
				t.Log(err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tt.expectedValues, r.k.Raw())
		})
	}
}

func TestRenderer_loadDataForFile(t *testing.T) {
	tests := map[string]struct {
		givenFileName  string
		givenValues    Values
		expectedValues Values
		expectErr      bool
	}{
		"GivenGlobalDefaults_WhenBuildingVarsForFile_ThenUseDefaults": {
			givenFileName: "README.md",
			givenValues: Values{
				":globals":  Values{"key": "variable"},
				"README.md": Values{"name": "readme"},
			},
			expectedValues: Values{"key": "variable", "name": "readme"},
		},
		"GivenNestedStructures_WhenBuildingVarsForFile_ThenDeepMergeAndOverride": {
			givenFileName: "README.md",
			givenValues: Values{
				":globals": Values{"key": "variable"},
				"README.md": Values{
					"key": Values{
						"nested": "key",
					}},
			},
			expectedValues: Values{"key": Values{"nested": "key"}},
		},
		"GivenSubDirectory_WhenBuildingVarsForFile_ThenMergeDirectoryValues": {
			givenFileName: "subdir/README.md",
			givenValues: Values{
				":globals": Values{"key": "variable"},
				"subdir": Values{
					"intermediate": "value",
				},
				"subdir/README.md": Values{
					"key": Values{
						"nested": "key",
					}},
			},
			expectedValues: Values{
				"key": Values{
					"nested": "key"},
				"intermediate": "value"},
		},
		"GivenMultipleSubDirectories_WhenBuildingVarsForFile_ThenMergeSubDirectoryValues": {
			givenFileName: "subdir1/subdir2/README.md",
			givenValues: Values{
				":globals": Values{"key": "variable"},
				"subdir1": Values{
					"intermediate1": "foo",
					"override":      "to-be-overridden",
				},
				"subdir1/subdir2": Values{
					"intermediate2": "bar",
					"override":      "inherited",
				},
				"subdir1/subdir2/README.md": Values{
					"key": Values{
						"nested": "key",
					}},
			},
			expectedValues: Values{
				"key": Values{
					"nested": "key"},
				"intermediate1": "foo",
				"intermediate2": "bar",
				"override":      "inherited",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				p: printer.DefaultPrinter,
				k: koanf.New("."),
			}
			require.NoError(t, r.k.Load(confmap.Provider(tt.givenValues, ""), nil))
			result, err := r.loadDataForFile(tt.givenFileName)
			if tt.expectErr {
				t.Log(err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tt.expectedValues, result)
		})
	}
}
