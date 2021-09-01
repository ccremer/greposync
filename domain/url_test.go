package domain

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gitUrlTests = map[string]struct {
	rawUrl                 string
	expectedRepoName       string
	expectedNamespace      string
	expectedFullName       string
	expectedRedactedString string
}{
	"GitHubURL": {
		rawUrl:                 "https://github.com/ccremer/greposync",
		expectedRepoName:       "greposync",
		expectedNamespace:      "ccremer",
		expectedFullName:       "github.com/ccremer/greposync",
		expectedRedactedString: "https://github.com/ccremer/greposync",
	},
	"GitLabURL": {
		rawUrl:                 "https://gitlab.com/gitlab-org/gitlab.git",
		expectedRepoName:       "gitlab",
		expectedNamespace:      "gitlab-org",
		expectedFullName:       "gitlab.com/gitlab-org/gitlab",
		expectedRedactedString: "https://gitlab.com/gitlab-org/gitlab.git",
	},
	"UserInfoURL": {
		rawUrl:                 "https://user:password@host.com:8443/namespace/repo.git",
		expectedRepoName:       "repo",
		expectedNamespace:      "namespace",
		expectedFullName:       "host.com:8443/namespace/repo",
		expectedRedactedString: "https://user:xxxxx@host.com:8443/namespace/repo.git",
	},
}

func TestGitUrl_GetRepositoryName(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitURL(*u)
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
			gitUrl := GitURL(*u)
			result := gitUrl.GetNamespace()
			assert.Equal(t, tt.expectedNamespace, result)
		})
	}
}

func TestGitUrl_Redacted(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitURL(*u)
			result := gitUrl.Redacted()
			assert.Equal(t, tt.expectedRedactedString, result)
		})
	}
}

func TestGitURL_GetFullName(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitURL(*u)
			result := gitUrl.GetFullName()
			assert.Equal(t, tt.expectedFullName, result)
		})
	}
}

func TestGitUrl_String(t *testing.T) {
	for name, tt := range gitUrlTests {
		t.Run(name, func(t *testing.T) {
			u, err := url.Parse(tt.rawUrl)
			require.NoError(t, err)
			gitUrl := GitURL(*u)
			result := gitUrl.String()
			assert.Equal(t, tt.rawUrl, result)
		})
	}
}
