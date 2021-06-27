package rendering

import (
	"bufio"
	"log"
	"os"
	"path"
	"text/template"
)

func RenderTemplate(repoDir string, data map[string]interface{}) error {
	funcs := funcMap()

	fileName := "README.tpl.md"
	targetFilePath := path.Join(repoDir, "README.md")
	tpl := template.Must(template.New(fileName).Option("missingkey=error").Funcs(funcs).ParseFiles("template/_helpers.tpl", "template/README.tpl.md"))

	f, err := os.Create(targetFilePath)
	CheckIfError(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tpl.Execute(w, data)
	CheckIfError(err)
	_ = w.Flush()
	return err
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
