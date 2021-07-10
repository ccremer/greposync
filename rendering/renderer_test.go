package rendering

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ccremer/greposync/cfg"
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
				"readme.md": "# example-repository\n\nEXAMPLE-REPOSITORY\ntest\n",
			},
		},
		"GivenTemplateFileInSubDir_WhenProcessing_ThenWriteFileToCorrectDir": {
			givenTemplate: "ci/pipeline.yml",
			expectedFileContents: map[string]string{
				"ci/pipeline.yml": "CommitBranch: \"\"\nCommitMessage: \"\"\nDefaultBranch: \"\"\nForcePush: false\nName: example-repository\nNamespace: \"\"\nSkipCommit: false\nSkipPush: false\nSkipReset: false\n",
			},
		},
	}
	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			r := &Renderer{
				p: printer.New(),
				cfg: &cfg.SyncConfig{
					Template: &cfg.TemplateConfig{RootDir: "testdata/template-1"},
					Git: &cfg.GitConfig{
						Dir:  s.TestGitDir,
						Name: "example-repository",
					},
				},
				k: s.K,
			}
			err := r.processTemplate(path.Join(r.cfg.Template.RootDir, tt.givenTemplate))
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

func (s *TemplateTestSuite) TestRenderer_ProcessTemplateDir() {
	r := &Renderer{
		p: printer.New(),
		cfg: &cfg.SyncConfig{
			Git: &cfg.GitConfig{
				Dir: s.TestGitDir,
			},
			Template: &cfg.TemplateConfig{
				RootDir: "testdata/template-1",
			},
		},
		k:              koanf.New("."),
		globalDefaults: s.K,
	}
	result := r.ProcessTemplateDir()()
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
	err := r.applyTemplate(targetPath, nil, Values{
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
	err := r.applyTemplate(targetPath, nil, Values{
		"Values": Values{
			"delete": true,
		}})
	s.Require().NoError(err)
	s.Assert().NoFileExists(targetPath)
}

func Test_sanitizeTargetPath(t *testing.T) {
	tests := map[string]struct {
		givenPath    string
		expectedPath string
	}{
		"GivenFileWithoutDir_WhenSanitizing_ThenReturnSamePath": {
			givenPath:    "fileName",
			expectedPath: "fileName",
		},
		"GivenFileInDir_WhenSanitizing_ThenReturnSamePath": {
			givenPath:    "dir/fileName",
			expectedPath: "dir/fileName",
		},
		"GivenFileWithTplExtension_WhenSanitizing_ThenReturnStripped": {
			givenPath:    "dir/fileName.tpl",
			expectedPath: "dir/fileName",
		},
		"GivenFileWithTplExtensionTwice_WhenSanitizing_ThenReturnStrippedOnce": {
			givenPath:    "fileName.tpl.tpl",
			expectedPath: "fileName.tpl",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := sanitizeTargetPath(tt.givenPath)
			assert.Equal(t, tt.expectedPath, result)
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
