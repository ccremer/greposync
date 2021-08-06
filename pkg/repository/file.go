package repository

import (
	"io/fs"
	"os"
	"path"

	"github.com/ccremer/greposync/core"
	"golang.org/x/sys/unix"
)

// GetConfig implements core.GitRepository.
func (g *Repository) GetConfig() core.GitRepositoryConfig {
	return g.coreConfig
}

// DeleteFile implements core.GitRepository.
func (g *Repository) DeleteFile(relativePath string) error {
	fileName := path.Join(g.Config.Dir, relativePath)
	if pathExists(fileName) {
		return os.Remove(fileName)
	}
	return nil
}

// EnsureFile implements core.GitRepository.
func (g *Repository) EnsureFile(targetFile, content string, fileMode fs.FileMode) error {
	// This allows us to create files with 777 permissions
	originalUmask := unix.Umask(0)
	defer unix.Umask(originalUmask)

	fileName := path.Join(g.Config.Dir, targetFile)

	if err := g.createParentDirs(fileName); err != nil {
		return err
	}

	// To ensure we can update the file permissions, as os.WriteFile does not change permissions.
	if err := g.DeleteFile(fileName); err != nil {
		return err
	}
	return os.WriteFile(fileName, []byte(content), fileMode)
}

func (g *Repository) createParentDirs(targetPath string) error {
	dir := path.Dir(targetPath)
	return os.MkdirAll(dir, 0775)
}

func pathExists(fileOrDir string) bool {
	if _, err := os.Stat(fileOrDir); err == nil {
		return true
	}
	return false
}
