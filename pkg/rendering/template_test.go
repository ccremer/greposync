package rendering

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplateService_evaluatePath(t *testing.T) {
	tests := map[string]struct {
		givenFilePath     string
		givenFileInfo     fs.FileInfo
		givenError        error
		expectedErrString string
		expectedResult    *core.Template
	}{
		"GivenNormalFile_ThenReturnWithSameMode": {
			givenFilePath: "template/filename",
			givenFileInfo: fakeFileInfo{mode: 0755},
			expectedResult: &core.Template{
				RelativePath: "filename",
				FileMode:     0755,
			},
		},
		"GivenFileInSubDir_ThenReturnRelativePathWithSubdir": {
			givenFilePath: "template/subdir/file",
			givenFileInfo: fakeFileInfo{mode: 0755},
			expectedResult: &core.Template{
				RelativePath: "subdir/file",
				FileMode:     0755,
			},
		},
		"GivenDir_ThenReturnRelativePathWithSubdir": {
			givenFilePath:  "template/subdir",
			givenFileInfo:  fakeFileInfo{dir: true},
			expectedResult: nil,
		},
		"GivenHelperFile_ThenReturnIgnore": {
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
			s := &GoTemplateService{config: &cfg.TemplateConfig{RootDir: "template"}}
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

func TestGoTemplateService_RenderTemplate(t *testing.T) {
	tests := map[string]struct {
		givenTemplate     core.Template
		expectedErrString string
		expectedTemplate  string
	}{
		"GivenExistingTemplate_ThenParseWithoutError": {
			givenTemplate:    core.Template{RelativePath: "template_1.tpl", FileMode: 0777},
			expectedTemplate: "variable\n\"json\"\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &GoTemplateService{
				config:    &cfg.TemplateConfig{RootDir: "testdata"},
				templates: map[string]*template.Template{},
			}
			err := s.parseTemplate(tt.givenTemplate)
			if tt.expectedErrString != "" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			fileName := "tmp_test_" + name
			defer os.Remove(fileName)

			renderErr := s.RenderTemplate(core.Output{
				Template:   tt.givenTemplate,
				TargetPath: fileName,
				Values: map[string]interface{}{
					"data": "json",
				},
			})
			require.NoError(t, renderErr)
			content, readErr := os.ReadFile(fileName)
			require.NoError(t, readErr)
			assert.Equal(t, tt.expectedTemplate, string(content))
			info, statErr := os.Stat(fileName)
			require.NoError(t, statErr)
			assert.Equal(t, tt.givenTemplate.FileMode, info.Mode())
		})
	}
}

func TestGoTemplateService_listAllTemplates(t *testing.T) {
	// We have to compute the file mode instead of hardcoding it, in CI/CD after cloning it might not be the same.
	fileInfo, err := os.Stat("testdata/template_1.tpl")
	require.NoError(t, err)
	s := &GoTemplateService{
		config: &cfg.TemplateConfig{RootDir: "testdata"},
	}
	result, err := s.listAllTemplates()
	require.NoError(t, err)
	assert.Equal(t, core.Template{
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
