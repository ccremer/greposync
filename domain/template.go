package domain

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

// Permissions is an alias for file permissions.
type Permissions fs.FileMode

// Template is a reference to a file that contains special syntax.
type Template struct {
	// RelativePath is the Path reference to where the template file is contained within the template root directory.
	RelativePath Path
	// FilePermissions defines what file permissions this template file has.
	// Rendered files should have the same permissions as template files.
	FilePermissions Permissions
}

// NewTemplate returns a new instance.
func NewTemplate(relPath Path, perms Permissions) *Template {
	return &Template{
		RelativePath:    relPath,
		FilePermissions: perms,
	}
}

// Render takes the given Values and returns a RenderResult from the given TemplateEngine.
func (t *Template) Render(values Values, engine TemplateEngine) (RenderResult, error) {
	content, err := engine.Execute(t, values)
	return content, err
}

// CleanPath returns a new Path with the first occurrence of ".tpl" in the base file name removed.
func (t *Template) CleanPath() Path {
	dirName := path.Dir(t.RelativePath.String())
	baseName := path.Base(t.RelativePath.String())
	newName := strings.Replace(baseName, ".tpl", "", 1)
	return NewPath(dirName, newName)
}

// AsValues returns the metadata as Values for rendering.
func (t *Template) AsValues() Values {
	return Values{
		"RelativePath": t.RelativePath.String(),
		"Permissions":  t.FilePermissions.Octal(),
	}
}

// FileMode converts Permissions to fs.FileMode.
func (p Permissions) FileMode() fs.FileMode {
	return fs.FileMode(p)
}

// Octal returns an octal permission representation (Linux)
func (p Permissions) Octal() string {
	return fmt.Sprintf("%#o", p.FileMode().Perm())
}
