package rendering

import (
	"errors"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/pkg/valuestore"
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
		valueStore     *valuestore.KoanfValueStore
	}
	Values     map[string]interface{}
	FileAction func(targetPath string, data Values) error
)

// NewRenderer returns a new instance of a renderer.
func NewRenderer(c *cfg.SyncConfig, globalDefaults *koanf.Koanf, cfg *cfg.Configuration) *Renderer {
	return &Renderer{
		p:              printer.New().SetLevel(printer.DefaultLevel).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.Git.Name),
		k:              koanf.New("."),
		globalDefaults: globalDefaults,
		cfg:            c,
		instance:       rendering.NewGoTemplateStore(cfg),
		valueStore:     valuestore.NewValueStore(),
	}
}

// RenderTemplateDir renders the templates parsed by ParseTemplateDir.
// Values are injected.
// Files are written to git target directory, although special Values may override that behaviour.
func (r *Renderer) RenderTemplateDir() pipeline.ActionFunc {
	return func() pipeline.Result {
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
	values, err := r.valueStore.FetchValuesForTemplate(tpl, &core.GitRepositoryProperties{
		URL:     core.FromURL(r.cfg.Git.Url),
		RootDir: r.cfg.Git.Dir,
	})
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
	if values["Values"].(core.Values)["delete"] == true {
		// files are deleted in a separate step
		return nil
	}
	gitCfg := core.GitRepositoryProperties{
		URL:     core.FromURL(r.cfg.Git.Url),
		RootDir: r.cfg.Git.Dir,
	}
	unmanaged, err := r.valueStore.FetchUnmanagedFlag(tpl, &gitCfg)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	if unmanaged {
		r.p.InfoF("Leaving file alone due to 'unmanaged' flag being set: %s", targetPath)
		return nil
	}

	newPath, err := r.valueStore.FetchTargetPath(tpl, &gitCfg)
	if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
		return err
	}
	if newPath != "" {
		r.p.DebugF("Redefining target path from '%s' to '%s", targetPath, newPath)
		targetPath = newPath
	}
	result, err := tpl.Render(values)
	if err != nil {
		return err
	}
	tpl.(*rendering.GoTemplate).RelativePath = targetPath
	g := repository.NewGitRepository(r.cfg.Git, r.cfg.PullRequest, nil)
	return g.EnsureFile(tpl.GetRelativePath(), result, tpl.GetFileMode())
}
