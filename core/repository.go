package core

import "io/fs"

// GitRepository is a domain entity enabling interaction with a local Git repository.
//counterfeiter:generate . GitRepository
type GitRepository interface {
	// FetchLabels returns a list of repository labels to be managed.
	FetchLabels() []Label
	// GetConfig returns the GitRepositoryConfig instance associated for this particular repository.
	GetConfig() GitRepositoryConfig

	// DeleteFile removes the given path from the Git repository relative to the root dir.
	// No error is returned if the file does not exist.
	DeleteFile(relativePath string) error
	// EnsureFile creates or updates a file in the Git repository.
	// targetPath is relative to the Git repository root dir.
	// content is the file content to write.
	// fileMode specifies the file permissions.
	EnsureFile(targetPath, content string, fileMode fs.FileMode) error
}
