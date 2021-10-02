package domain

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGolden_MetadataValues(t *testing.T) {
	u, err := url.Parse("https://github.com/ccremer/greposync")
	require.NoError(t, err)
	gitUrl := FromURL(u)

	fileName := "testdata/golden/metadata.tpl"

	repository := NewGitRepository(gitUrl, NewPath("repos", gitUrl.GetFullName()))
	repository.CommitBranch = "greposync-update"
	repository.DefaultBranch = "master"
	ctx := RenderContext{
		Repository: repository,
	}
	template := NewTemplate(NewPath(fileName), Permissions(0644))

	file, err := os.Create(fileName)
	require.NoError(t, err)
	defer file.Close()
	for k, _ := range ctx.Repository.AsValues() {
		placeholder := fmt.Sprintf("{{ .%s.%s.%s }}", MetadataValueKey, RepositoryValueKey, k)
		_, _ = fmt.Fprintf(file, "{{`%s`}} = %s\n", placeholder, placeholder)
	}
	for k, _ := range template.AsValues() {
		placeholder := fmt.Sprintf("{{ .%s.%s.%s }}", MetadataValueKey, TemplateValueKey, k)
		_, _ = fmt.Fprintf(file, "{{`%s`}} = %s\n", placeholder, placeholder)
	}

	values := ctx.enrichWithMetadata(Values{}, template)
	engine := DummyEngine{templatePath: fileName}
	result, err := engine.Execute(template, values)
	require.NoError(t, err)
	err = result.WriteToFile(NewPath("testdata", "golden", "metadata.txt"), template.FilePermissions)
	require.NoError(t, err)
}
