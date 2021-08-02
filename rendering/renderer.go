package rendering

import (
	"path"
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
)

type (
	Renderer struct {
		p              printer.Printer
		cfg            *cfg.SyncConfig
		k              *koanf.Koanf
		globalDefaults *koanf.Koanf
		instance       *rendering.GoTemplateStore
	}
	Values     map[string]interface{}
	FileAction func(targetPath string, data Values) error
)

var (
	SyncConfigFileName = ".sync.yml"
)

// NewRenderer returns a new instance of a renderer.
func NewRenderer(c *cfg.SyncConfig, globalDefaults *koanf.Koanf) *Renderer {
	return &Renderer{
		p:              printer.New().SetLevel(printer.DefaultLevel).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.Git.Name),
		k:              koanf.New("."),
		globalDefaults: globalDefaults,
		cfg:            c,
		instance:       rendering.NewGoTemplateStore(c.Template),
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

		templates, err := r.instance.FetchTemplates()
		if err != nil {
			return pipeline.Result{Err: err}
		}
		for _, tpl := range templates {
			if err = r.processTemplate(tpl); err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) processTemplate(tpl core.Template) error {
	values, err := r.loadDataForFile(tpl.GetRelativePath())
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Values":   values,
		"Metadata": r.ConstructMetadata(),
	}

	targetPath := tpl.GetRelativePath()
	return r.applyTemplate(targetPath, tpl, data)
}

func (r *Renderer) applyTemplate(targetPath string, tpl core.Template, values core.Values) error {
	if values["Values"].(Values)["delete"] == true {
		return r.deleteFileIfExists(targetPath)
	}
	if values["Values"].(Values)["unmanaged"] == true {
		r.p.InfoF("Leaving file alone due to 'unmanaged' flag being set: %s", targetPath)
		return nil
	}
	if newTarget := values["Values"].(Values)["targetPath"]; newTarget != nil && newTarget != "" {
		newPath := newTarget.(string)
		if strings.HasSuffix(newPath, "/") {
			newPath = path.Clean(path.Join(newPath, path.Base(targetPath)))
		} else {
			newPath = path.Clean(path.Join(newPath))
		}
		r.p.DebugF("Redefining target path from '%s' to '%s", targetPath, newPath)
		targetPath = newPath
	}
	result, err := tpl.Render(values)
	if err != nil {
		return err
	}
	tpl.(*rendering.GoTemplate).RelativePath = targetPath
	g := repository.NewGitRepository(r.cfg.Git, nil)
	return g.EnsureFile(tpl.GetRelativePath(), result, tpl.GetFileMode())
}
