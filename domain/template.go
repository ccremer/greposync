package domain

import (
	"io/fs"
	"os"
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

func (t *Template) RenderToFile(values Values, absolutePath Path, renderer TemplateRenderer) error {
	content, err := renderer.Render(values, t)
	if err != nil {
		return err
	}
	return os.WriteFile(absolutePath.String(), []byte(content), t.FilePermissions.FileMode())
}

func (p Permissions) FileMode() fs.FileMode {
	return fs.FileMode(p)
}
