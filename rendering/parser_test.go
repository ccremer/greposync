package rendering

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_GivenTemplateDirectory_WhenListingAllFiles_IgnoreHelpers(t *testing.T) {
	r := &Parser{}
	root := "testdata/template-1"
	results, err := r.listAllFiles(root)
	require.NoError(t, err)
	assert.Equal(t, []string{
		path.Join(root, "ci", "pipeline.yml"),
		path.Join(root, "readme.tpl.md"),
	}, results)
}
