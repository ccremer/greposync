package domain

import "os"

// RenderResult represents the string value after rendering from a Template.
type RenderResult string

// WriteToFile writes the content to the given Path with given Permissions.
// Otherwise, an error is returned.
func (r RenderResult) WriteToFile(path Path, permissions Permissions) error {
	return os.WriteFile(path.String(), []byte(r), permissions.FileMode())
}
