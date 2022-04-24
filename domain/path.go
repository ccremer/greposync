package domain

import (
	"os"
	"path"
	"path/filepath"
)

// Path is a Value object identifying a file path.
type Path string

// NewPath constructs a new Path joined by the given elements.
// Paths are joined with path.Join.
func NewPath(elems ...string) Path {
	return Path(path.Join(elems...))
}

// NewFilePath constructs a new Path joined by the given elements.
// Paths are joined with filepath.Join.
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

// Join takes this Path as root and makes a new Path with given elements.
func (p Path) Join(elems ...Path) Path {
	var strElems = make([]string, len(elems)+1)
	strElems[0] = p.String()
	for i := range elems {
		strElems[i+1] = elems[i].String()
	}
	return NewFilePath(strElems...)
}

// Delete removes the path (and possibly all children if it's a directory), ignoring any errors.
// If you need error handling, use os.RemoveAll directly.
func (p Path) Delete() {
	_ = os.RemoveAll(p.String())
}

// String returns a string representation of itself.
func (p Path) String() string {
	return string(p)
}

// IsInSlice returns true if p is in the given slice, false otherwise.
func (p Path) IsInSlice(paths []Path) bool {
	for i := 0; i < len(paths); i++ {
		if paths[i] == p {
			return true
		}
	}
	return false
}
