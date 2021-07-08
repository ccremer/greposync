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

// ProcessTemplateDir searches for template files in the configure dir, renders the template with injected data and writes them to git target directory.
func (r *Renderer) ProcessTemplateDir() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := r.loadVariables(path.Join(r.cfg.Git.Dir, ".sync.yml"))
		if err != nil {
			return pipeline.Result{Err: err}
		}

		files, err := r.listAllFiles(path.Clean(r.cfg.Template.RootDir))
		if err != nil {
			return pipeline.Result{Err: err}
		}
		for _, file := range files {
			if err = r.processTemplate(file); err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) processTemplate(originalTemplatePath string) error {
	relativePath, err := filepath.Rel(r.cfg.Template.RootDir, sanitizeTargetPath(originalTemplatePath))
	if err != nil {
		return err
	}
	originalFileName := path.Base(originalTemplatePath)

	templates := []string{originalTemplatePath}
	helperFilePath := path.Join(r.cfg.Template.RootDir, helperFileName)
	if r.fileExists(helperFilePath) {
		templates = append(templates, helperFilePath)
	}
	// Read template and helpers
	tpl, err := template.
		New(originalFileName).
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
	targetPath := path.Join(r.cfg.Git.Dir, relativePath)
	return r.writeFile(targetPath, tpl, data)
}

func (r *Renderer) writeFile(targetPath string, tpl *template.Template, data Values) error {
	r.p.InfoF("Writing file: %s", path.Base(targetPath))
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

func (r *Renderer) listAllFiles(root string) (files []string, err error) {
	err = filepath.Walk(root,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileName := path.Base(file)
			// Don't add helper file or directories
			if fileName != helperFileName && r.fileExists(file) {
				files = append(files, file)
			}
			return nil
		})
	return files, err
}

func sanitizeTargetPath(targetPath string) string {
	dirName := path.Dir(targetPath)
	baseName := path.Base(targetPath)
	newName := strings.Replace(baseName, ".tpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}
