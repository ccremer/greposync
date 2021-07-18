package initialize

import (
	_ "embed"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

var (
	//go:embed _helpers.tpl
	helperTpl []byte
	//go:embed README.md.tpl
	readmeTpl []byte

	//go:embed greposync.yml
	grepoSyncYml []byte
	//go:embed config_defaults.yml
	configDefaultsYml []byte
	//go:embed managed_repos.yml
	managedReposYml []byte

	configFiles = map[string][]byte{
		"greposync.yml":       grepoSyncYml,
		"config_defaults.yml": configDefaultsYml,
		"managed_repos.yml":   managedReposYml,
	}
	templateFiles = map[string][]byte{
		"template/_helpers.tpl": helperTpl,
		"template/README.md":    readmeTpl,
	}
)

// CreateMainConfigFiles creates the main configuration files.
// Each pre existing file is skipped.
func CreateMainConfigFiles() pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: writeFiles(configFiles)}
	}
}

// CreateTemplateFiles creates the example files in the template directory.
// The dir has to exist.
func CreateTemplateFiles() pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: writeFiles(templateFiles)}
	}
}

func writeFiles(files map[string][]byte) error {
	for file, content := range files {
		err := writeFile(file, content)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(file string, content []byte) error {
	if !fileExists(file) {
		return os.WriteFile(file, content, 0644)
	}
	return nil
}

// CreateTemplateDir creates the template directory if it doesn't exist.
func CreateTemplateDir() pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: createDir("template")}
	}
}

func fileExists(path string) bool {
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		return true
	}
	return false
}

func dirExists(path string) bool {
	if f, err := os.Stat(path); err == nil && f.IsDir() {
		return true
	}
	return false
}

func createDir(path string) error {
	if !dirExists(path) {
		return os.Mkdir(path, 0775)
	}
	return nil
}
