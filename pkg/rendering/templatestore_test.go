package rendering

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplateStore_ImplementsInterface(t *testing.T) {
	assert.Implements(t, (*core.TemplateStore)(nil), new(GoTemplateStore))
}

func TestGoTemplateStore_evaluatePath(t *testing.T) {
	tests := map[string]struct {
		givenFilePath     string
		givenFileInfo     fs.FileInfo
		givenError        error
		expectedErrString string
		expectedResult    *GoTemplate
	}{
		"GivenNormalFile_ThenReturnWithSameMode": {
			givenFilePath: "template/filename",
			givenFileInfo: fakeFileInfo{mode: 0755},
			expectedResult: &GoTemplate{
				RelativePath: "filename",
				FileMode:     0755,
			},
		},
		"GivenFileInSubDir_ThenReturnRelativePathWithSubdir": {
			givenFilePath: "template/subdir/file",
			givenFileInfo: fakeFileInfo{mode: 0755},
			expectedResult: &GoTemplate{
				RelativePath: "subdir/file",
				FileMode:     0755,
			},
		},
		"GivenDir_ThenReturnRelativePathWithSubdir": {
			givenFilePath:  "template/subdir",
			givenFileInfo:  fakeFileInfo{dir: true},
			expectedResult: nil,
		},
		"GivenHelperFile_ThenIgnore": {
			givenFilePath:  "template/_helpers.tpl",
			givenFileInfo:  fakeFileInfo{},
			expectedResult: nil,
		},
		"GivenIOError_ThenReturnError": {
			givenFilePath:     "template/subdir",
			givenFileInfo:     fakeFileInfo{},
			givenError:        errors.New("io error"),
			expectedErrString: "io error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &GoTemplateStore{config: &cfg.Configuration{Template: &cfg.TemplateConfig{RootDir: "template"}}}
			result, err := s.evaluatePath(tt.givenFilePath, tt.givenFileInfo, tt.givenError)
			if tt.expectedErrString != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrString)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGoTemplateStore_listAllTemplates(t *testing.T) {
	// We have to compute the file mode instead of hardcoding it, in CI/CD after cloning it might not be the same.
	fileInfo, err := os.Stat("testdata/template_1.tpl")
	require.NoError(t, err)
	s := &GoTemplateStore{
		config: &cfg.Configuration{Template: &cfg.TemplateConfig{RootDir: "testdata"}},
	}
	result, err := s.listAllTemplates()
	require.NoError(t, err)
	assert.Equal(t, &GoTemplate{
		RelativePath: "template_1.tpl",
		FileMode:     fileInfo.Mode(),
	}, result[0])
	assert.Len(t, result, 1)
}

type fakeFileInfo struct {
	dir  bool
	mode fs.FileMode
}

func (f fakeFileInfo) Name() string {
	panic("implement me")
}

func (f fakeFileInfo) Size() int64 {
	panic("implement me")
}

func (f fakeFileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f fakeFileInfo) ModTime() time.Time {
	panic("implement me")
}

func (f fakeFileInfo) IsDir() bool {
	return f.dir
}

func (f fakeFileInfo) Sys() interface{} {
	panic("implement me")
}
