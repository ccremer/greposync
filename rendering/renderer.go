package rendering

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ccremer/git-repo-sync/printer"
	"github.com/knadh/koanf"
)

type (
	Service struct {
		p   printer.Printer
		cfg Config
		k   *koanf.Koanf
	}
	Config struct {
		TemplateDir    string
		TargetRootDir  string
		PrTemplatePath string
		RepoName       string
		ConfigDefaults *koanf.Koanf
	}
	Data map[string]interface{}
)

const (
	errorOnMissingKey = "missingkey=error"
)

var (
	templateFunctions = funcMap()
)

func NewService(c Config) *Service {
	return &Service{
		p:   printer.New().SetLevel(printer.LevelDebug).MapColorToLevel(printer.Magenta, printer.LevelInfo).SetName(c.RepoName),
		k:   koanf.New("."),
		cfg: c,
	}
}

func (s *Service) ProcessTemplates() {
	s.loadVariables(path.Join(s.cfg.TargetRootDir, ".sync.yml"))

	files, err := filepath.Glob(path.Join(s.cfg.TemplateDir, "*"))
	s.p.CheckIfError(err)
	for _, file := range files {
		fileName := path.Base(file)
		if strings.HasPrefix(fileName, "_") && strings.HasSuffix(fileName, ".tpl") {
			// File is a helper file
			continue
		}
		s.processTemplate(file)
	}
}

func (s *Service) processTemplate(fullPath string) {
	relativePath, err := filepath.Rel(s.cfg.TemplateDir, fullPath)
	s.p.CheckIfError(err)
	targetPath := path.Join(s.cfg.TargetRootDir, relativePath)
	targetPath = sanitizeTargetPath(targetPath)
	fileName := path.Base(relativePath)

	// Read template and helpers
	tpl, err := template.
		New(fileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(fullPath)
	s.p.CheckIfError(err)
	tpl, _ = tpl.ParseGlob(path.Join(s.cfg.TemplateDir, "_*.tpl"))

	data := map[string]interface{}{
		"Values":   s.loadDataForFile(relativePath),
		"Metadata": s.getMetadata(targetPath),
	}

	// Write target file
	s.p.InfoF("Writing file: %s", path.Base(targetPath))
	f, err := os.Create(targetPath)
	printer.CheckIfError(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	err = tpl.Execute(w, data)
	printer.CheckIfError(err)
	_ = w.Flush()
}

func (s *Service) getMetadata(targetPath string) Data {
	return map[string]interface{}{
		"Path":     targetPath,
		"FileName": path.Base(targetPath),
		"RepoName": path.Base(s.cfg.TargetRootDir),
	}
}

func sanitizeTargetPath(targetPath string) string {
	dirName := path.Dir(targetPath)
	baseName := path.Base(targetPath)
	newName := strings.Replace(baseName, ".gotmpl", "", 1)
	return path.Clean(path.Join(dirName, newName))
}
