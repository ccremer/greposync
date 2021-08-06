package core

import "errors"

//go:generate rm -r corefakes
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//CoreService is a representation of a core feature or process.
type CoreService interface {
	// RunPipeline executes the main business logic of this core service.
	// It returns an error if the core service deems the process to have failed.
	RunPipeline() error
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

// TemplateStore is a service responsible for fetching templates.
type TemplateStore interface {
	// FetchTemplates retrieves the templates or an error if one failed.
	FetchTemplates() ([]Template, error)
}

// ValueStore is a service centered around configuration values fetching and configuring templates.
type ValueStore interface {
	// FetchValuesForTemplate retrieves the Values for the given template.
	FetchValuesForTemplate(template Template, config *GitRepositoryConfig) (Values, error)
	// FetchUnmanagedFlag returns true if the given template should not be rendered.
	// The implementation may return ErrKeyNotFound if the flag is undefined, as the boolean 'false' is ambiguous.
	FetchUnmanagedFlag(template Template, config *GitRepositoryConfig) (bool, error)
	// FetchTargetPath returns an alternative output path for the given template relative to the Git repository.
	// An empty string indicates that there is no alternative path configured.
	FetchTargetPath(template Template, config *GitRepositoryConfig) (string, error)
	// FetchFilesToDelete returns a slice of paths that should be deleted in the Git repository.
	// The paths are relative to the Git directory.
	FetchFilesToDelete(config *GitRepositoryConfig) ([]string, error)
}

// ErrKeyNotFound is an error that indicates that a particular key was not found in the ValueStore.
var ErrKeyNotFound = errors.New("key not found")
