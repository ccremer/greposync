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
	}, s.K, &cfg.Configuration{
		Template: &cfg.TemplateConfig{
			RootDir: "testdata/template-1",
		},
	})
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
