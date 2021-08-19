package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var pathTests = map[string]struct {
	path   Path
	isFile bool
	isDir  bool
}{
	"ExistingDirectory": {
		path:  "testdata",
		isDir: true,
	},
	"ExistingFile": {
		path:   Path(filepath.Join("testdata", "path_test")),
		isFile: true,
	},
	"NonExistingPath": {
		path: Path(filepath.Join("testdata", "non-existing")),
	},
}

func TestPath_Exists(t *testing.T) {
	for name, tt := range pathTests {
		t.Run(name, func(t *testing.T) {
			result := tt.path.Exists()
			shouldExist := tt.isFile || tt.isDir
			assert.Equal(t, shouldExist, result)
		})
	}
}

func TestPath_FileExists(t *testing.T) {
	for name, tt := range pathTests {
		t.Run(name, func(t *testing.T) {
			result := tt.path.FileExists()
			assert.Equal(t, tt.isFile, result)
		})
	}
}

func TestPath_DirExists(t *testing.T) {
	for name, tt := range pathTests {
		t.Run(name, func(t *testing.T) {
			result := tt.path.DirExists()
			assert.Equal(t, tt.isDir, result)
		})
	}
}

func TestPath_Delete(t *testing.T) {
	t.Run("DeleteFile", func(t *testing.T) {
		f, err := os.CreateTemp("testdata", "PATH_")
		require.NoError(t, err)
		assert.FileExists(t, f.Name())
		assert.NoError(t, f.Close())
		NewFilePath(f.Name()).Delete()
		assert.NoFileExists(t, f.Name())
	})
	t.Run("DeleteDir", func(t *testing.T) {
		dir, err := os.MkdirTemp("testdata", "PATH_DIR_")
		require.NoError(t, err)
		f, err := os.CreateTemp(dir, "PATH_")
		require.NoError(t, err)
		assert.FileExists(t, f.Name())
		assert.NoError(t, f.Close())
		NewFilePath(dir).Delete()
		assert.NoFileExists(t, f.Name())
		assert.NoDirExists(t, dir)
	})
}

func TestPath_Join(t *testing.T) {
	tests := map[string]struct {
		givenPath      Path
		givenElems     []Path
		expectedResult Path
	}{
		"GivenEmptyPath_WhenNilElems_ThenExpectEmpty": {
			givenPath:      "",
			givenElems:     nil,
			expectedResult: "",
		},
		"GivenEmptyPath_WhenEmptyElems_ThenExpectEmpty": {
			givenPath:      "",
			givenElems:     []Path{},
			expectedResult: "",
		},
		"GivenEmptyPath_WhenSingleElems_ThenExpectSingle": {
			givenPath:      "",
			givenElems:     []Path{"single"},
			expectedResult: "single",
		},
		"GivenEmptyPath_WhenMultipleElems_ThenExpectCombined": {
			givenPath:      "",
			givenElems:     []Path{"multi", "path"},
			expectedResult: "multi/path",
		},
		"GivenSinglePath_WhenMultipleElems_ThenExpectWithRoot": {
			givenPath:      "root",
			givenElems:     []Path{"multi", "path"},
			expectedResult: "root/multi/path",
		},
		"GivenSubDirPath_WhenMultipleElems_ThenExpectWithRoot": {
			givenPath:      "root/subdir",
			givenElems:     []Path{"multi", "path"},
			expectedResult: "root/subdir/multi/path",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.givenPath.Join(tt.givenElems...)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
