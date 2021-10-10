package domain

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
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

	repositoryValues := ctx.Repository.AsValues().Keys()
	sort.Strings(repositoryValues)
	printLines(MetadataValueKey, RepositoryValueKey, repositoryValues, file)

	templateValues := template.AsValues().Keys()
	sort.Strings(templateValues)
	printLines(MetadataValueKey, TemplateValueKey, templateValues, file)

	values := ctx.enrichWithMetadata(Values{}, template)
	engine := DummyEngine{templatePath: fileName}
	result, err := engine.Execute(template, values)
	require.NoError(t, err)
	err = result.WriteToFile(NewPath("testdata", "golden", "metadata.txt"), template.FilePermissions)
	require.NoError(t, err)
}

func printLines(firstKey, secondKey string, thirdKeys []string, file io.Writer) {
	for _, thirdKey := range thirdKeys {
		placeholder := fmt.Sprintf("{{ .%s.%s.%s }}", firstKey, secondKey, thirdKey)
		_, _ = fmt.Fprintf(file, "{{`%s`}} = %s\n", placeholder, placeholder)
	}
}
