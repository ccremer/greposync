package pullrequest

import (
	"errors"
	"net/url"
	"testing"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/corefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestService_fetchPrTemplate(t *testing.T) {
	tests := map[string]struct {
		expectTemplate bool
	}{
		"GivenNoPrTemplate_ThenExpectEmptyTemplate": {
			expectTemplate: false,
		},
		"GivenExistingPrTemplate_WhenSuccessful_ThenExpectTemplate": {
			expectTemplate: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			template := newFakeTemplate("", nil)
			templateStore := &corefakes.FakeTemplateStore{
				FetchPullRequestTemplateStub: func() (core.Template, error) {
					if tt.expectTemplate {
						return template, nil
					}
					return nil, nil
				},
			}
			s := &PullRequestService{
				templateStore: templateStore,
			}
			ctx := &pipelineContext{
				repo: &corefakes.FakeGitRepository{GetConfigStub: func() core.GitRepositoryConfig {
					return core.GitRepositoryConfig{
						URL: toUrl(t, "https://github.com/ccremer/greposync"),
					}
				}},
			}
			result := s.fetchPrTemplate(ctx)
			assert.Equal(t, 1, templateStore.FetchPullRequestTemplateCallCount())
			if tt.expectTemplate {
				assert.Equal(t, template, ctx.template)
			} else {
				assert.Nil(t, ctx.template)
			}
			require.NoError(t, result)
		})
	}
}

func toUrl(t *testing.T, raw string) *core.GitURL {
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return core.FromURL(u)
}

func TestPullRequestService_renderTemplate(t *testing.T) {
	tests := map[string]struct {
		givenTemplate *corefakes.FakeTemplate
		expectedBody  string
		expectedError error
	}{
		"GivenNoTemplate_WhenRender_ThenExpectEmptyBody": {
			givenTemplate: nil,
			expectedBody:  "",
		},
		"GivenTemplate_WhenRender_ThenExpectABody": {
			givenTemplate: newFakeTemplate("rendered", nil),
			expectedBody:  "rendered",
		},
		"GivenTemplate_WhenRenderFails_ThenExpectError": {
			givenTemplate: newFakeTemplate("", errors.New("failed")),
			expectedBody:  "rendered",
			expectedError: errors.New("failed"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &PullRequestService{}
			ctx := &pipelineContext{
				template: tt.givenTemplate,
				repo: &corefakes.FakeGitRepository{GetConfigStub: func() core.GitRepositoryConfig {
					return core.GitRepositoryConfig{URL: toUrl(t, "github.com/ccremer/greposync")}
				}},
			}
			err := s.renderTemplate(ctx)
			if tt.expectedBody != "" {
				assert.Equal(t, 1, tt.givenTemplate.RenderCallCount())
			}
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBody, ctx.body)
		})
	}
}

func newFakeTemplate(returnString string, returnErr error) *corefakes.FakeTemplate {
	return &corefakes.FakeTemplate{RenderStub: func(values core.Values) (string, error) {
		return returnString, returnErr
	}}
}
