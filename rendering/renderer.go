package rendering

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
)

type (
	Renderer struct {
		p              printer.Printer
		cfg            *cfg.SyncConfig
		k              *koanf.Koanf
		globalDefaults *koanf.Koanf
		parser         *Parser
	}
	Values     map[string]interface{}
	FileAction func(targetPath string, data Values) error
)

var (
	templateFunctions  = funcMap()
	SyncConfigFileName = ".sync.yml"
)

// NewRenderer returns a new instance of a renderer.
func NewRenderer(c *cfg.SyncConfig, globalDefaults *koanf.Koanf, parser *Parser) *Renderer {
	return &Renderer{
		p:              printer.New().SetLevel(printer.DefaultLevel).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.Git.Name),
		k:              koanf.New("."),
		globalDefaults: globalDefaults,
		cfg:            c,
		parser:         parser,
	}
}

// RenderTemplateDir renders the templates parsed by ParseTemplateDir.
// Values from SyncConfigFileName are injected.
// Files are written to git target directory, although special Values may override that behaviour.
func (r *Renderer) RenderTemplateDir() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := r.loadVariables(path.Join(r.cfg.Git.Dir, SyncConfigFileName))
		if err != nil {
			return pipeline.Result{Err: err}
		}

		for file, tpl := range r.parser.templates {
			if err = r.processTemplate(file, tpl); err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) processTemplate(originalTemplatePath string, tpl *template.Template) error {
	relativePath, err := filepath.Rel(r.cfg.Template.RootDir, cleanTargetPath(originalTemplatePath))
	if err != nil {
		return err
	}

	values, err := r.loadDataForFile(relativePath)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Values":   values,
		"Metadata": r.ConstructMetadata(),
	}

	targetPath := path.Join(r.cfg.Git.Dir, relativePath)
	return r.applyTemplate(targetPath, tpl, data)
}

func (r *Renderer) applyTemplate(targetPath string, tpl *template.Template, values Values) error {
	if values["Values"].(Values)["delete"] == true {
		if fileExists(targetPath) {
			r.p.DebugF("Deleting file due to 'delete' flag being set: %s", targetPath)
			return os.Remove(targetPath)
		}
		return nil
	}
	if values["Values"].(Values)["unmanaged"] == true {
		r.p.DebugF("Leaving file alone due to 'unmanaged' flag being set: %s", targetPath)
		return nil
	}
	if newTarget := values["Values"].(Values)["targetPath"]; newTarget != nil && newTarget != "" {
		newPath := newTarget.(string)
		if strings.HasSuffix(newPath, string(filepath.Separator)) {
			newPath = path.Clean(path.Join(r.cfg.Git.Dir, newPath, path.Base(targetPath)))
		} else {
			newPath = path.Clean(path.Join(r.cfg.Git.Dir, newPath))
		}
		r.p.DebugF("Redefining target path from '%s' to '%s", targetPath, newPath)
		targetPath = newPath
	}
	return r.writeFile(targetPath, tpl, values)
}

func (r *Renderer) writeFile(targetPath string, tpl *template.Template, data Values) error {
	r.p.InfoF("Writing file from template: %s", path.Base(targetPath))
	dir := path.Dir(targetPath)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}
	return w.Flush()
}

func cleanTargetPath(targetPath string) string {
	dirName := path.Dir(targetPath)
	baseName := path.Base(targetPath)
	newName := strings.Replace(baseName, ".tpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}
