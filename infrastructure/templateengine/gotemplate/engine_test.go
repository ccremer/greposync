package gotemplate

import (
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoTemplateEngine_ExecuteString(t *testing.T) {
	expected := "title"

	engine := NewEngine()
	values := domain.Values{
		"Values": domain.Values{
			"title": expected,
		},
	}
	templateString := "{{ .Values.title }}"

	result, err := engine.ExecuteString(templateString, values)
	require.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
