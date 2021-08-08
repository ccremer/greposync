package repository

import (
	"io/fs"
	"os"
	"path"

	"github.com/ccremer/greposync/core"
	"golang.org/x/sys/unix"
)

// GetConfig implements core.GitRepository.
func (s *Repository) GetConfig() core.GitRepositoryProperties {
	return s.coreConfig
}

// DeleteFile implements core.GitRepository.
func (s *Repository) DeleteFile(relativePath string) error {
	fileName := path.Join(s.GitConfig.Dir, relativePath)
	if pathExists(fileName) {
		return os.Remove(fileName)
	}
	return nil
}

// EnsureFile implements core.GitRepository.
func (s *Repository) EnsureFile(targetFile, content string, fileMode fs.FileMode) error {
	// This allows us to create files with 777 permissions
	originalUmask := unix.Umask(0)
	defer unix.Umask(originalUmask)

	fileName := path.Join(s.GitConfig.Dir, targetFile)

	if err := s.createParentDirs(fileName); err != nil {
		return err
	}

	// To ensure we can update the file permissions, as os.WriteFile does not change permissions.
	if err := s.DeleteFile(fileName); err != nil {
		return err
	}
	return os.WriteFile(fileName, []byte(content), fileMode)
}

func (s *Repository) createParentDirs(targetPath string) error {
	dir := path.Dir(targetPath)
	return os.MkdirAll(dir, 0775)
}

func pathExists(fileOrDir string) bool {
	if _, err := os.Stat(fileOrDir); err == nil {
		return true
	}
	return false
}
