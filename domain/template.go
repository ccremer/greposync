package domain

import (
	"io/fs"
)

type Permissions fs.FileMode

type Template struct {
	RelativePath    Path
	FilePermissions Permissions
}

func NewTemplate(relPath Path, perms Permissions) *Template {
	return &Template{
		RelativePath:    relPath,
		FilePermissions: perms,
	}
}

func (t *Template) Render(values Values, engine TemplateEngine) (string, error) {
	content, err := engine.Execute(t, values)
	return content, err
}

func (p Permissions) FileMode() fs.FileMode {
	return fs.FileMode(p)
}
