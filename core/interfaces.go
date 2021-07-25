package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//CoreService is a representation of a core feature or process.
type CoreService interface {
	// RunPipeline executes the main business logic of this core service.
	// It returns an error if the core service deems the process to have failed.
	RunPipeline() error
}

// ManagedRepoProvider is a core service that is responsible for providing services for managing Git repositories.
//counterfeiter:generate . ManagedRepoProvider
type ManagedRepoProvider interface {
	// LoadManagedRepositories will load the managed repositories from a config store and returns an array of GitRepositoryFacade for each Git repository.
	LoadManagedRepositories() ([]GitRepositoryFacade, error)
	// GetSupportedGitHostingProviders returns a map of supported GitHostingFacade.
	GetSupportedGitHostingProviders() map[GitHostingProvider]GitHostingFacade
}

// GitRepositoryFacade is a core service enabling interaction with a local Git repository.
//counterfeiter:generate . GitRepositoryFacade
type GitRepositoryFacade interface {
	// GetLabels returns a list of repository labels to be managed.
	GetLabels() []GitRepositoryLabel
	// GetConfig returns the GitRepositoryConfig instance associated for this particular repository.
	GetConfig() GitRepositoryConfig
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
	CreateOrUpdateLabelsForRepo(url *GitUrl, labels []GitRepositoryLabel) error
	/*
		DeleteLabelsForRepo deletes the given labels from the given remote repository.
		An error may be returned on first deletion failure, but non-existing labels are not errors.
	*/
	DeleteLabelsForRepo(url *GitUrl, labels []GitRepositoryLabel) error
}

// GitRepositoryLabel is a label that is attached to a remote Git repository on a supported Git hosting provider.
// The implementation may contain additional provider-specific properties.
//counterfeiter:generate . GitRepositoryLabel
type GitRepositoryLabel interface {
	// GetName returns the name of the label.
	// Used for internal identification, may not the be label name that is actually in the remote repository.
	GetName() string
	// IsBoundForDeletion returns true if the label is bound for removal from a remote repository.
	IsBoundForDeletion() bool
}
