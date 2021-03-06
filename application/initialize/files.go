package initialize

import (
	"context"
	_ "embed"
	"os"
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
)

// createMainConfigFiles creates the main configuration files.
// Each pre-existing file is skipped.
func (c *Command) createMainConfigFiles(_ context.Context) error {
	return writeFiles(c.configFiles)
}

// createTemplateFiles creates the example files in the template directory.
// The dir has to exist.
func (c *Command) createTemplateFiles(_ context.Context) error {
	return writeFiles(c.templateFiles)
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

// createTemplateDir creates the template directory if it doesn't exist.
func (c *Command) createTemplateDir(_ context.Context) error {
	return createDir(c.TemplateDir)
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
