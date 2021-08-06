package repository

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/ccremer/greposync/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitRepository_DeleteFile(t *testing.T) {
	tests := map[string]struct {
		givenPath  string
		createFile bool
	}{
		"GivenExistingFile_ThenDeleteIt": {
			givenPath:  "testdata/tmp_" + t.Name(),
			createFile: true,
		},
		"GivenNonExistingFile_ThenIgnoreError": {
			givenPath: "testdata/tmp_" + t.Name(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.createFile {
				require.NoError(t, os.WriteFile(tt.givenPath, []byte("test"), 0644))
			}
			g := &Repository{
				Config: &cfg.GitConfig{Dir: "testdata"},
			}
			err := g.DeleteFile(path.Base(tt.givenPath))
			require.NoError(t, err)
			assert.NoFileExists(t, tt.givenPath)
		})
	}
}

func TestGitRepository_EnsureFile(t *testing.T) {
	tests := map[string]struct {
		givenPath    string
		expectedMode fs.FileMode
	}{
		"GivenFilePermissions_ThenTemporarilyDisableUmask": {
			givenPath:    "tmp1_" + t.Name(),
			expectedMode: 0777,
		},
		"GivenNormalFilePermissions_ThenCreateFile": {
			givenPath:    "tmp2_" + t.Name(),
			expectedMode: 0644,
		},
		"GivenMissingDirectories_ThenCreateParentDirs": {
			givenPath:    "subdir/tmp3_" + t.Name(),
			expectedMode: 0644,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			g := &Repository{
				Config: &cfg.GitConfig{Dir: "testdata"},
			}
			fullPath := path.Join("testdata", tt.givenPath)
			err := g.EnsureFile(tt.givenPath, "test", tt.expectedMode)
			require.NoError(t, err)
			assert.FileExists(t, fullPath)
			info, err := os.Stat(fullPath)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMode, info.Mode())
			require.NoError(t, os.Remove(fullPath))
		})
	}
}
