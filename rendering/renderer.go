package rendering

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/printer"
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

func (r *Renderer) ProcessTemplates() {
	r.loadVariables(path.Join(r.cfg.Git.Dir, ".sync.yml"))

	files, err := filepath.Glob(path.Join(r.cfg.Template.RootDir, "*"))
	r.p.CheckIfError(err)
	for _, file := range files {
		fileName := path.Base(file)
		if fileName == helperFileName {
			// File is a helper file
			continue
		}
		r.processTemplate(file)
	}
}

func (r *Renderer) processTemplate(templateFullPath string) {
	relativePath, err := filepath.Rel(r.cfg.Template.RootDir, templateFullPath)
	r.p.CheckIfError(err)
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
	r.p.CheckIfError(err)

	data := map[string]interface{}{
		"Values":   r.loadDataForFile(relativePath),
		"Metadata": r.ConstructMetadata(),
	}

	// Write target file
	r.p.InfoF("Writing file: %s", path.Base(targetPath))
	f, err := os.Create(targetPath)
	printer.CheckIfError(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tpl.Execute(w, data)
	printer.CheckIfError(err)
	_ = w.Flush()
}

func sanitizeTargetPath(targetPath string) string {
	dirName := path.Dir(targetPath)
	baseName := path.Base(targetPath)
	newName := strings.Replace(baseName, ".gotmpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}
