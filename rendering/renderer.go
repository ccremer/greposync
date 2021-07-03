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

func NewRenderer(c *cfg.SyncConfig, globalDefaults *koanf.Koanf) *Renderer {
	return &Renderer{
		p:              printer.New().SetLevel(printer.LevelDebug).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.Git.Name),
		k:              koanf.New("."),
		globalDefaults: globalDefaults,
		cfg:            c,
	}
}

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
	if r.FileExists(helperFilePath) {
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
