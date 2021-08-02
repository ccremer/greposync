package rendering

import (
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TemplateTestSuite struct {
	suite.Suite
	TestGitDir    string
	SeedSourceDir string
	SeedTargetDir string
	K             *koanf.Koanf
}

func TestRenderer(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

func (s *TemplateTestSuite) dirExists(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

func (s *TemplateTestSuite) SetupTest() {
	s.TestGitDir = "testdata/template-1-test"
	s.SeedSourceDir = "testdata/template-2"
	s.SeedTargetDir = "testdata/template-2-test"
	assert.NoError(s.T(), os.RemoveAll(s.TestGitDir))
	assert.NoError(s.T(), os.RemoveAll(s.SeedTargetDir))
	assert.NoError(s.T(), os.MkdirAll(s.TestGitDir, 0755))
	assert.NoError(s.T(), os.MkdirAll(s.SeedTargetDir, 0755))
	values := Values{
		"readme.md": Values{
			"custom": "test",
		}}
	k := koanf.New(".")
	s.Require().NoError(k.Load(confmap.Provider(values, ""), nil))
	s.copyFiles()
	s.K = k
}

func (s *TemplateTestSuite) TearDownTest() {
	if !s.T().Failed() {
		s.Assert().NoError(os.RemoveAll(s.TestGitDir))
		s.Assert().NoError(os.RemoveAll(s.SeedTargetDir))
	}
}

func (s *TemplateTestSuite) TestProcessTemplate() {
	tests := map[string]struct {
		givenTemplate        string
		givenValues          Values
		expectedFileContents map[string]string
		expectErr            bool
	}{
		"GivenTemplateFile_WhenProcessing_ThenWriteFile": {
			givenTemplate: "readme.tpl.md",
			expectedFileContents: map[string]string{
				"readme.md": "EXAMPLE-REPOSITORY",
			},
		},
		"GivenTemplateFileInSubDir_WhenProcessing_ThenWriteFileToCorrectDir": {
			givenTemplate: "ci/pipeline.yml",
			expectedFileContents: map[string]string{
				"ci/pipeline.yml": "EXAMPLE-REPOSITORY",
			},
		},
	}
	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			u, err := url.Parse("https://github.com/example/example-repository")
			require.NoError(t, err)
			r := NewRenderer(&cfg.SyncConfig{
				Template: &cfg.TemplateConfig{RootDir: "testdata/template-1"},
				Git: &cfg.GitConfig{
					Dir:  s.TestGitDir,
					Name: "example-repository",
					Url:  u,
				},
			}, s.K)
			require.NoError(t, err)
			err = r.processTemplate(&rendering.GoTemplate{
				RelativePath: tt.givenTemplate,
				FileMode:     0644,
			})
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			for fileName, expectedContent := range tt.expectedFileContents {
				content, readErr := os.ReadFile(path.Join(s.TestGitDir, fileName))
				require.NoError(t, readErr)
				assert.Equal(t, expectedContent, string(content))
			}
		})
	}
}

func (s *TemplateTestSuite) TestRenderer_RenderTemplateDir() {
	u, err := url.Parse("https://github.com/example/irrelevant")
	require.NoError(s.T(), err)
	r := NewRenderer(&cfg.SyncConfig{
		Git: &cfg.GitConfig{
			Dir: s.TestGitDir,
			Url: u,
		},
		Template: &cfg.TemplateConfig{
			RootDir: "testdata/template-1",
		},
	}, s.K)
	err = s.K.Load(confmap.Provider(map[string]interface{}{
		"readme.md": core.Values{
			"custom": "value",
		},
	}, ""), nil)
	s.Require().NoError(err)
	result := r.RenderTemplateDir()()
	s.Require().NoError(result.Err)
	s.Assert().NoFileExists(path.Join(s.TestGitDir, "_helpers.tpl"))
	s.Assert().FileExists(path.Join(s.TestGitDir, "readme.md"))
	s.Assert().FileExists(path.Join(s.TestGitDir, "ci", "pipeline.yml"))
}

func (s *TemplateTestSuite) TestRenderer_GivenUnmanagedFlag_WhenApplyingTemplate_ThenLeaveFileAlone() {
	r := &Renderer{
		p: printer.New(),
	}
	targetPath := path.Join(s.SeedTargetDir, "readme.md")
	err := r.applyTemplate(targetPath, &rendering.GoTemplate{}, core.Values{
		"Values": Values{
			"unmanaged": true,
		}})
	s.Require().NoError(err)
	s.Assert().FileExists(targetPath)
}

func (s *TemplateTestSuite) TestRenderer_GivenDeleteFlag_WhenApplyingTemplate_ThenRemoveTargetFileInstead() {
	r := &Renderer{
		p: printer.New(),
	}
	targetPath := path.Join(s.SeedTargetDir, "readme.md")
	err := r.applyTemplate(targetPath, &rendering.GoTemplate{}, core.Values{
		"Values": Values{
			"delete": true,
		}})
	s.Require().NoError(err)
	s.Assert().NoFileExists(targetPath)
}

func (s *TemplateTestSuite) TestRenderer_GivenTargetPath() {
	tests := map[string]struct {
		givenTargetDir              string
		expectedEffectiveTargetFile string
	}{
		"GivenTargetDir_WhenApplyingTemplate_ThenChangeDirectoryButNotFileName": {
			givenTargetDir:              "dir/",
			expectedEffectiveTargetFile: path.Join(s.SeedTargetDir, "dir", "readme.md"),
		},
		"GivenTargetPath_WhenApplyingTemplate_ThenChangeFileName": {
			givenTargetDir:              "dir/newFile.ext",
			expectedEffectiveTargetFile: path.Join(s.SeedTargetDir, "dir", "newFile.ext"),
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			u, err := url.Parse("https://github.com/example/irrelevant")
			require.NoError(t, err)
			r := NewRenderer(&cfg.SyncConfig{
				Template: &cfg.TemplateConfig{RootDir: s.SeedSourceDir},
				Git: &cfg.GitConfig{
					Dir: s.SeedTargetDir,
					Url: u,
				},
			}, koanf.New("."))
			targetPath := "readme.md"
			templates, err := r.instance.FetchTemplates()
			require.NoError(t, err)
			err = r.applyTemplate(targetPath, templates[0], core.Values{
				"Values": Values{
					"targetPath": tt.givenTargetDir,
				}})
			require.NoError(t, err)
			assert.FileExists(t, tt.expectedEffectiveTargetFile)
		})
	}
}

func (s *TemplateTestSuite) copyFiles() {
	files, err := filepath.Glob(s.SeedSourceDir + "/*")
	s.Require().NoError(err)
	for _, file := range files {
		source, openErr := os.Open(file)
		s.Require().NoError(openErr)

		target, tgtErr := os.Create(path.Join(s.SeedTargetDir, path.Base(source.Name())))
		s.Require().NoError(tgtErr)
		_, copyErr := io.Copy(target, source)
		s.Require().NoError(copyErr)
		s.Require().NoError(target.Close())
		s.Require().NoError(source.Close())
	}
}
