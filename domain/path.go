package domain

import (
	"os"
	"path"
	"path/filepath"
)

type Path string

func NewPath(elems ...string) Path {
	return Path(path.Join(elems...))
}

func NewFilePath(elems ...string) Path {
	return Path(filepath.Join(elems...))
}

// Exists returns true if the path exists in the local file system.
func (p Path) Exists() bool {
	if _, err := os.Stat(p.String()); err == nil {
		return true
	}
	return false
}

// FileExists returns true if the path exists in the local file system and is a file.
func (p Path) FileExists() bool {
	if info, err := os.Stat(p.String()); err == nil {
		return !info.IsDir()
	}
	return false
}

// DirExists returns true if the path exists in the local file system and is a directory.
func (p Path) DirExists() bool {
	if info, err := os.Stat(p.String()); err == nil {
		return info.IsDir()
	}
	return false
}

// Delete removes the path (and possibly all children if it's a directory), ignoring any errors.
// If you need error handling, use os.RemoveAll directly.
func (p Path) Delete() {
	_ = os.RemoveAll(p.String())
}

func (p Path) String() string {
	return string(p)
}
