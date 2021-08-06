package core

import "io/fs"

// GitRepository is a domain entity enabling interaction with a local Git repository.
//counterfeiter:generate . GitRepository
type GitRepository interface {
	// GetLabels returns a list of repository labels to be managed.
	// This method does not return the labels that are actually in the remote Git hosting service, but the ones configured locally.
	GetLabels() []Label
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

// GitRepositoryConfig holds all the relevant Git properties.
type GitRepositoryConfig struct {
	// URL is the repository location on the remote hosting provider.
	URL *GitURL
	// RootDir is the local root path to the Git repository.
	RootDir string
}

// GitRepositoryStore is a core service that is responsible for providing services for managing Git repositories.
//counterfeiter:generate . GitRepositoryStore
type GitRepositoryStore interface {
	// FetchGitRepositories will load the managed repositories from a config store and returns an array of GitRepository for each Git repository.
	// A non-nil empty slice is returned if there are none existing.
	FetchGitRepositories() ([]GitRepository, error)
}

// Label is attached to a remote Git repository on a supported Git hosting provider.
// The implementation may contain additional provider-specific properties.
//counterfeiter:generate . Label
type Label interface {
	// IsInactive returns true if the label is bound for removal from a remote repository.
	IsInactive() bool
	// GetName returns the label name.
	GetName() string

	// Delete removes the label from the remote repository.
	Delete() (bool, error)
	// Ensure creates the label in the remote repository if it doesn't exist.
	// If the label already exists, it will be updated if the properties are different.
	Ensure() (bool, error)
}
