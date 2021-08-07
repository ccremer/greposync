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

	// FetchPullRequest queries the remote Git hosting for an existing PullRequest.
	// If no PR with matching branches is found, then nil is returned without error.
	FetchPullRequest() (PullRequest, error)
	// NewPullRequest creates a new instance with default properties.
	NewPullRequest() PullRequest

	// EnsurePullRequest creates the given PullRequest if it doesn't exist.
	// The PullRequest is updated if it needs updating, otherwise left unchanged without error.
	EnsurePullRequest(pr PullRequest) error
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
