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
	}
	Values map[string]interface{}
)

const (
	errorOnMissingKey = "missingkey=error"
	helperFileName    = "_helpers.tpl"
)

var (
	templateFunctions = funcMap()
)

// NewRenderer returns a new instance of a renderer.
func NewRenderer(c *cfg.SyncConfig, globalDefaults *koanf.Koanf) *Renderer {
	return &Renderer{
		p:              printer.New().SetLevel(printer.DefaultLevel).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.Git.Name),
		k:              koanf.New("."),
		globalDefaults: globalDefaults,
		cfg:            c,
	}
}

// RenderPrTemplate renders the PR template.
// If the BodyTemplate config is a path to an existing file, it will use the file and overwrite the config with the rendered result.
// IF the BodyTemplate config is a string, it will overwrite it with a rendered and data-injected string.
// If the BodyTemplate config is empty, it will use the CommitMessage.
func (r *Renderer) RenderPrTemplate() pipeline.ActionFunc {
	return func() pipeline.Result {
		t := r.cfg.PullRequest.BodyTemplate
		if t == "" {
			r.p.InfoF("No PullRequest template defined")
			r.cfg.PullRequest.BodyTemplate = r.cfg.Git.CommitMessage
		}

		data := Values{"Metadata": r.ConstructMetadata()}
		if r.fileExists(t) {
			if str, err := r.RenderTemplateFile(data, t); err != nil {
				return pipeline.Result{Err: err}
			} else {
				r.cfg.PullRequest.BodyTemplate = str
			}
		} else {
			if str, err := r.RenderString(data, t); err != nil {
				return pipeline.Result{Err: err}
			} else {
				r.cfg.PullRequest.BodyTemplate = str
			}
		}
		return pipeline.Result{}
	}
}

// ProcessTemplates searches for template files in the configure dir, renders the template with injected data and writes them to git target directory.
func (r *Renderer) ProcessTemplates() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := r.loadVariables(path.Join(r.cfg.Git.Dir, ".sync.yml"))
		if err != nil {
			return pipeline.Result{Err: err}
		}

		files, err := filepath.Glob(path.Join(r.cfg.Template.RootDir, "*"))
		if err != nil {
			return pipeline.Result{Err: err}
		}
		for _, file := range files {
			fileName := path.Base(file)
			if fileName == helperFileName {
				// File is a helper file
				continue
			}
			if err = r.processTemplate(file); err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) processTemplate(templateFullPath string) error {
	relativePath, err := filepath.Rel(r.cfg.Template.RootDir, templateFullPath)
	if err != nil {
		return err
	}
	targetPath := path.Join(r.cfg.Git.Dir, relativePath)
	targetPath = sanitizeTargetPath(targetPath)
	fileName := path.Base(relativePath)

	templates := []string{templateFullPath}
	helperFilePath := path.Join(r.cfg.Template.RootDir, helperFileName)
	if r.fileExists(helperFilePath) {
		templates = append(templates, helperFilePath)
	}
	// Read template and helpers
	tpl, err := template.
		New(fileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(templates...)
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

	// Write target file
	return r.writeFile(targetPath, tpl, data)
}

func (r *Renderer) writeFile(targetPath string, tpl *template.Template, data Values) error {
	r.p.InfoF("Writing file: %s", path.Base(targetPath))
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

func sanitizeTargetPath(targetPath string) string {
	dirName := path.Dir(targetPath)
	baseName := path.Base(targetPath)
	newName := strings.Replace(baseName, ".gotmpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}
