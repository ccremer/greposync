package rendering

import (
	"bytes"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/ccremer/greposync/core"
)

// GoTemplate implements core.Template.
type GoTemplate struct {
	RelativePath string
	FileMode     fs.FileMode
	template     *template.Template
}

// GetRelativePath implements core.Template.
func (g *GoTemplate) GetRelativePath() string {
	dirName := path.Dir(g.RelativePath)
	baseName := path.Base(g.RelativePath)
	newName := strings.Replace(baseName, ".tpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}

// GetFileMode implements core.Template.
func (g *GoTemplate) GetFileMode() fs.FileMode {
	return g.FileMode
}

// Render implements core.Template.
func (g *GoTemplate) Render(values core.Values) (string, error) {
	buf := &bytes.Buffer{}
	err := g.template.Execute(buf, values)
	return buf.String(), err
}
