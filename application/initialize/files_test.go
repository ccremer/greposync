package initialize

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/ccremer/greposync/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			for file := range tt.givenFiles {
				assert.FileExists(t, file)
				require.NoError(t, os.Remove(file))
			}
		})
	}
}

func Test_CommentHelper(t *testing.T) {
	tests := map[string]struct {
		givenValues     domain.Values
		expectedContent string
	}{
		"GivenEmptyComments": {
			givenValues: domain.Values{
				"text":   "",
				"open":   "",
				"closed": "",
				"prefix": "",
			},
			expectedContent: "# my-repository\n\ndescription\n",
		},
		"GivenEmptyCommentsWithEnclosures": {
			givenValues: domain.Values{
				"text":        "",
				"open":        "<!--",
				"closed":      "-->",
				"prefix":      "",
				"name":        "my-repository",
				"description": "My awesome, greposync managed repository",
			},
			expectedContent: "# my-repository\n\ndescription\n",
		},
		"GivenSingleComments": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync",
				"open":        "",
				"closed":      "",
				"prefix":      "",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "This file is managed by greposync\n# my-repository\n\ndescription\n",
		},
		"GivenMultilineComments": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync\nDo not edit",
				"open":        "",
				"closed":      "",
				"prefix":      "",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "This file is managed by greposync\nDo not edit\n# my-repository\n\ndescription\n",
		},
		"GivenEnclosedComment": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync",
				"open":        "<!--",
				"closed":      "-->",
				"prefix":      "",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "<!--\nThis file is managed by greposync\n-->\n# my-repository\n\ndescription\n",
		},
		"GivenEnclosedCommentWithPrefixes": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync",
				"open":        "<!--",
				"closed":      "-->",
				"prefix":      "# ",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "<!--\n# This file is managed by greposync\n-->\n# my-repository\n\ndescription\n",
		},
		"GivenEnclosedMultilineComment": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync\nDo not edit",
				"open":        "<!--",
				"closed":      "-->",
				"prefix":      "",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "<!--\nThis file is managed by greposync\nDo not edit\n-->\n# my-repository\n\ndescription\n",
		},
		"GivenEnclosedMultilineCommentWithPrefixes": {
			givenValues: domain.Values{
				"text":        "This file is managed by greposync\nDo not edit",
				"open":        "<!--",
				"closed":      "-->",
				"prefix":      "# ",
				"name":        "my-repository",
				"description": "description",
			},
			expectedContent: "<!--\n# This file is managed by greposync\n# Do not edit\n-->\n# my-repository\n\ndescription\n",
		},
	}
	readme := "README.md.tpl"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tpl, err := template.New(readme).Funcs(sprig.TxtFuncMap()).ParseFiles("_helpers.tpl", readme)
			require.NoError(t, err)
			buf := bytes.NewBuffer([]byte{})
			err = tpl.Execute(buf, map[string]domain.Values{
				"Values": {
					"comment":     tt.givenValues,
					"name":        "my-repository",
					"description": "description",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedContent, buf.String())
		})
	}
}
