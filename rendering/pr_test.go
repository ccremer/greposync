package rendering

import (
	"testing"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_RenderString(t *testing.T) {
	tests := map[string]struct {
		givenTemplate  string
		givenValues    Values
		expectedString string
		expectErr      bool
	}{
		"GivenEmptyTemplate_WhenRendering_ThenReturnSameText": {
			givenTemplate:  "Update from control repository",
			givenValues:    Values{},
			expectedString: "Update from control repository",
		},
		"GivenGoTemplate_WhenRendering_ThenReturnRenderedText": {
			givenTemplate:  "Update from {{ .Values.name }}",
			givenValues:    Values{"Values": Values{"name": "control repository"}},
			expectedString: "Update from control repository",
		},
		"GivenMissingData_WhenRendering_ThenReturnError": {
			givenTemplate: "Update from {{ .Values.name }}",
			givenValues:   Values{"Values": Values{}},
			expectErr:     true,
		},
		"GivenInvalidTemplate_WhenRendering_ThenReturnError": {
			givenTemplate: "Invalid go template syntax {{ .Values.name }",
			givenValues:   Values{"Values": Values{}},
			expectErr:     true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				p: printer.DefaultPrinter,
			}
			result, err := r.renderString(tt.givenValues, tt.givenTemplate)
			if tt.expectErr {
				t.Log(err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedString, result)
		})
	}
}

func TestRenderer_RenderTemplateFile(t *testing.T) {
	tests := map[string]struct {
		givenFilePath  string
		givenValues    Values
		expectedString string
		expectErr      bool
	}{
		"GivenFileTemplate_WhenRendering_ThenReturnFileContent": {
			givenFilePath: "testdata/pr-template-1.md",
			givenValues: Values{"Metadata": Values{
				"Repository": Values{
					"Name": "template",
				}}},
			expectedString: "This is a multiline PR template that is called 'template'.\n\nIt has multiple lines\n",
		},
		"GivenInvalidTemplate_WhenRendering_ThenReturnError": {
			givenFilePath: "testdata/pr-template-1.md",
			givenValues:   Values{"not": Values{}},
			expectErr:     true,
		},
		"GivenNonExistingTemplate_WhenRendering_ThenReturnError": {
			givenFilePath: "not-existing",
			givenValues:   Values{},
			expectErr:     true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				p: printer.DefaultPrinter,
			}
			result, err := r.renderTemplateFile(tt.givenValues, tt.givenFilePath)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedString, result)
		})
	}
}

func TestRenderer_RenderPrTemplate(t *testing.T) {
	tests := map[string]struct {
		givenTemplate    string
		expectedTemplate string
		expectErr        bool
	}{
		"GivenFileNameInTemplate_WhenRendering_ThenRenderFromFile": {
			givenTemplate:    "testdata/pr-template-1.md",
			expectedTemplate: "This is a multiline PR template that is called 'example-repository'.\n\nIt has multiple lines\n",
		},
		"GivenEmptyTemplate_WhenRendering_ThenRenderCommitMessage": {
			expectedTemplate: "CommitMessage",
		},
		"GivenInlineTemplate_WhenRendering_ThenRenderFromInline": {
			givenTemplate:    "{{ .Metadata.Repository.Name }}",
			expectedTemplate: "example-repository",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &Renderer{
				p: printer.New(),
				cfg: &cfg.SyncConfig{
					Git:         &cfg.GitConfig{CommitMessage: "CommitMessage", Name: "example-repository"},
					PullRequest: &cfg.PullRequestConfig{BodyTemplate: tt.givenTemplate},
				},
				k: koanf.New("."),
			}
			result := r.RenderPrTemplate()()
			if tt.expectErr {
				require.Error(t, result.Err)
				return
			}
			require.NoError(t, result.Err)
			assert.Equal(t, tt.expectedTemplate, r.cfg.PullRequest.BodyTemplate)
		})
	}
}
