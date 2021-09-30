package domain

// GitRepositoryStore provides methods to interact with GitRepository on the local filesystem.
// Most methods described follow the corresponding Git operations.
//
// In Domain-Driven Design language, the term `Store` corresponds to `Repository`, but to avoid name clash it was named `Store`.
type GitRepositoryStore interface {
	// FetchGitRepositories loads a list of GitRepository from a configuration set.
	// Returns an empty list on first error.
	FetchGitRepositories() ([]*GitRepository, error)

	// Clone will download the given GitRepository to local filesystem.
	// The location is specified in GitRepository.RootDir.
	Clone(repository *GitRepository) error
	// Checkout checks out the GitRepository.CommitBranch.
	Checkout(repository *GitRepository) error
	// Fetch retrieves the objects and refs from remote.
	Fetch(repository *GitRepository) error
	// Reset current HEAD to GitRepository.CommitBranch.
	Reset(repository *GitRepository) error
	// Pull integrates objects from remote.
	Pull(repository *GitRepository) error

	// Add stages all files in GitRepository.RootDir.
	Add(repository *GitRepository) error
	// Commit records changes in the repository.
	Commit(repository *GitRepository, options CommitOptions) error
	// Diff returns a `patch`-compatible diff between HEAD and previous commit as string.
	// The diff may be empty.
	Diff(repository *GitRepository) (string, error)

	// Push updates remote refs.
	Push(repository *GitRepository, options PushOptions) error
}

// CommitOptions contains settings to influence the GitRepositoryStore.Commit action.
type CommitOptions struct {
	// Message contains the commit message.
	Message string
	// Amend will edit the last commit instead of creating a new one.
	Amend bool
}

// PushOptions contains settings to influence the GitRepositoryStore.Push action.
type PushOptions struct {
	// Force overwrites the remote state when pushing.
	Force bool
}
