package domain

// PullRequestStore provides methods to interact with PullRequest on a Git hosting service.
//
// In Domain-Driven Design language, the term `Store` corresponds to `Repository`, but to avoid name clash it was named `Store`.
type PullRequestStore interface {
	// FindMatchingPullRequest returns the PullRequest that has the same branch as GitRepository.CommitBranch.
	// If not found, it returns nil without error.
	FindMatchingPullRequest(repository *GitRepository) (*PullRequest, error)

	// EnsurePullRequest creates or updates the GitRepository.PullRequest in the repository.
	//
	//  * This operation does not alter any properties of existing labels.
	//  * Existing labels are left untouched, but any extraneous labels are removed.
	//  * Title and Body are updated.
	//  * Existing Commit and Base branches are left untouched.
	//
	// The first error encountered aborts the operation.
	EnsurePullRequest(repository *GitRepository) error
}
