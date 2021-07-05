package rendering

import (
	"os"
	"path"
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
	TestGitDir string
	K          *koanf.Koanf
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
	s.TestGitDir = "testdata/template_test"
	assert.NoError(s.T(), os.RemoveAll(s.TestGitDir))
	assert.NoError(s.T(), os.Mkdir(s.TestGitDir, 0755))
	values := Values{
		"readme.md": Values{
			"custom": "test",
		}}
	k := koanf.New(".")
	s.Require().NoError(k.Load(confmap.Provider(values, ""), nil))
	s.K = k
}

func (s *TemplateTestSuite) TearDownTest() {
	if !s.T().Failed() {
		assert.NoError(s.T(), os.RemoveAll(s.TestGitDir))
	}
}

func (s *TemplateTestSuite) TestProcessTemplate() {
	tests := map[string]struct {
		givenTemplate        string
		expectedFileContents map[string]string
		expectErr            bool
	}{
		"GivenTemplateFile_WhenProcessing_ThenWriteFile": {
			givenTemplate: "readme.gotmpl.md",
			expectedFileContents: map[string]string{
				"readme.md": "# example-repository\n\nEXAMPLE-REPOSITORY\ntest\n",
			},
		},
		"GivenTemplateFileInSubDir_WhenProcessing_ThenWriteFileToCorrectDir": {
			givenTemplate: "ci/pipeline.yml",
			expectedFileContents: map[string]string{
				"ci/pipeline.yml": "CommitBranch: \"\"\nCommitMessage: \"\"\nCreatePR: false\nDefaultBranch: \"\"\nForcePush: false\nName: example-repository\nNamespace: \"\"\nSkipCommit: false\nSkipPush: false\nSkipReset: false\n",
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
		"GivenFileWithGotmplExtension_WhenSanitizing_ThenReturnStripped": {
			givenPath:    "dir/fileName.gotmpl",
			expectedPath: "dir/fileName",
		},
		"GivenFileWithGotmplExtensionTwice_WhenSanitizing_ThenReturnStrippedOnce": {
			givenPath:    "fileName.gotmpl.gotmpl",
			expectedPath: "fileName.gotmpl",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := sanitizeTargetPath(tt.givenPath)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}
