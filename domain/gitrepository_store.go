package domain

type GitRepositoryStore interface {
	FetchGitRepositories() ([]*GitRepository, error)

	Clone(repository *GitRepository) error
	Checkout(repository *GitRepository) error
	Fetch(repository *GitRepository) error
	Reset(repository *GitRepository) error
	Pull(repository *GitRepository) error

	Add(repository *GitRepository) error
	Commit(repository *GitRepository, options CommitOptions) error
	Diff(repository *GitRepository) (string, error)
}

type CommitOptions struct {
	Message string
	Amend   bool
}
