package rendering

import (
	"bufio"
	"os"
	"path"
	"text/template"

	"github.com/ccremer/git-repo-sync/printer"
)

func RenderTemplate(repoDir string, data map[string]interface{}) error {
	funcs := funcMap()

	fileName := "README.tpl.md"
	targetFilePath := path.Join(repoDir, "README.md")
	tpl := template.Must(template.New(fileName).Option("missingkey=error").Funcs(funcs).ParseFiles("template/_helpers.tpl", "template/README.tpl.md"))

	f, err := os.Create(targetFilePath)
	printer.CheckIfError(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tpl.Execute(w, data)
	printer.CheckIfError(err)
	_ = w.Flush()
	return err
}
