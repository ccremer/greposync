package rendering

import (
	"bufio"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func RenderTemplate(repoDir string, data map[string]interface{}) error {
	funcs := sprig.TxtFuncMap()

	fileName := "README.tpl.md"
	targetFilePath := path.Join("repos", repoDir, "README.md")
	tpl := template.Must(template.New(fileName).Funcs(funcs).ParseFiles("template/_helpers.tpl", "template/README.tpl.md"))

	f, _ := os.Create(targetFilePath)
	defer f.Close()

	w := bufio.NewWriter(f)

	err := tpl.Execute(w, data)
	_ = w.Flush()
	return err
}
