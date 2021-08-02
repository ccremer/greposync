package core

import "io/fs"

// Values is a key-value construct with arbitrary hierarchy.
type Values map[string]interface{}

// Template is a representation of a single template file.
type Template interface {
	// GetRelativePath returns the path to a template file relative to the template root directory.
	// The path is delimited with a forward slash ("/") and not OS-specific.
	GetRelativePath() string
	// GetFileMode returns the mode of the template file.
	GetFileMode() fs.FileMode
	// Render takes the given Values and returns a string that essentially is the content of a file.
	Render(values Values) (string, error)
}
