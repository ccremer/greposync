package core

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gitUrlTests = map[string]struct {
	rawUrl            string
	expectedRepoName  string
	expectedNamespace string
}{
	"GitHubURL": {
		rawUrl:            "https://github.com/ccremer/greposync",
		expectedRepoName:  "greposync",
		expectedNamespace: "ccremer",
	},
	"GitLabURL": {
		rawUrl:            "https://gitlab.com/gitlab-org/gitlab.git",
		expectedRepoName:  "gitlab",
		expectedNamespace: "gitlab-org",
	},
}

func TestGitUrl_GetRepositoryName(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitUrl(*u)
			result := gitUrl.GetRepositoryName()
			assert.Equal(t, tt.expectedRepoName, result)
		})
	}
}

func TestGitUrl_GetNamespace(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitUrl(*u)
			result := gitUrl.GetNamespace()
			assert.Equal(t, tt.expectedNamespace, result)
		})
	}
}
