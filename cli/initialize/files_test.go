package initialize

import (
	"os"
	"testing"

	"github.com/ccremer/greposync/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func Test_ParseGreposyncYaml(t *testing.T) {
	bytes, err := os.ReadFile("greposync.yml")
	require.NoError(t, err)
	config := &cfg.Configuration{}
	err = yaml.Unmarshal(bytes, config)
	assert.NoError(t, err)
	assert.Equal(t, []string{"greposync"}, config.PullRequest.Labels)
	assert.False(t, config.PullRequest.Create)
	assert.Equal(t, "Update from greposync", config.PullRequest.Subject)
}

func Test_createDir(t *testing.T) {
	tests := map[string]struct {
		givenDir  string
		expectErr bool
	}{
		"GivenNonExistingDirectory_WhenCreating_ThenCreateDirectory": {
			givenDir: "testdir",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := createDir(tt.givenDir)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.DirExists(t, tt.givenDir)
			defer require.NoError(t, os.Remove(tt.givenDir))
		})
	}
}

func Test_writeFile(t *testing.T) {
	tests := map[string]struct {
		givenFilename string
		givenContent  string
		expectErr     bool
	}{
		"GivenNonExistingFile_WhenWriting_ThenCreateFileWithContent": {
			givenFilename: "test_file",
			givenContent:  "test content",
		},
		"GivenInvalidFileName_WhenWriting_ThenExpectError": {
			givenFilename: "invalid/",
			expectErr:     true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := writeFile(tt.givenFilename, []byte(tt.givenContent))
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			result, readErr := os.ReadFile(tt.givenFilename)
			defer require.NoError(t, os.Remove(tt.givenFilename))
			require.NoError(t, readErr)
			assert.Equal(t, []byte(tt.givenContent), result)
		})
	}
}

func Test_writeFiles(t *testing.T) {
	tests := map[string]struct {
		givenFiles map[string][]byte
		expectErr  bool
	}{
		"GivenNonExistingFiles_WhenWriting_ThenCreateFile": {
			givenFiles: map[string][]byte{
				"test_file": []byte("content"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := writeFiles(tt.givenFiles)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for file, _ := range tt.givenFiles {
				assert.FileExists(t, file)
				require.NoError(t, os.Remove(file))
			}
		})
	}
}
