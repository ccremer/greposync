package domain

import "os"

type RenderResult string

func (r RenderResult) WriteToFile(path Path, permissions Permissions) error {
	return os.WriteFile(path.String(), []byte(r), permissions.FileMode())
}
