package core

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
	// GetSupportedGitHostingProviders returns a map of supported GitHostingFacade.
	// A non-nil empty map is returned if there are none.
	GetSupportedGitHostingProviders() map[GitHostingProvider]GitHostingFacade
}

// GitHostingFacade is a core service providing interaction with remote Git hosting services.
//counterfeiter:generate . GitHostingFacade
type GitHostingFacade interface {
	// Initialize will initialize this service as required by the underlying provider.
	// An error shall be returned when it's not safe to continue interacting with the provider.
	Initialize() error
	/*
		CreateOrUpdateLabelsForRepo updates the repository labels or creates them if not existing.
		It is up to the underlying instance to figure out which labels and their config.
		An error may be returned if any operation failed.
		This method only mutates the given labels without deletions.
		Other labels are left ignored.
	*/
	CreateOrUpdateLabelsForRepo(url *GitURL, labels []Label) error
	/*
		DeleteLabelsForRepo deletes the given labels from the given remote repository.
		An error may be returned on first deletion failure, but non-existing labels are not errors.
	*/
	DeleteLabelsForRepo(url *GitURL, labels []Label) error
}

// Label is attached to a remote Git repository on a supported Git hosting provider.
// The implementation may contain additional provider-specific properties.
//counterfeiter:generate . Label
type Label interface {
	// IsBoundForDeletion returns true if the label is bound for removal from a remote repository.
	IsBoundForDeletion() bool

	// Delete removes the label from the remote repository.
	//Delete() error
	// Ensure creates the label in the remote repository if it doesn't exist.
	// If the label already exists, it will be updated if the properties are different.
	//Ensure() error
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
	FetchUnmanagedFlag(template Template, config *GitRepositoryConfig) bool
	// FetchTargetPath returns an alternative output path for the given template relative to the Git repository.
	// An empty string indicates that there is no alternative path configured.
	FetchTargetPath(template Template, config *GitRepositoryConfig) string
	// FetchFilesToDelete returns a slice of paths that should be deleted in the Git repository.
	// The paths are relative to the Git directory.
	FetchFilesToDelete(config *GitRepositoryConfig) ([]string, error)
}
