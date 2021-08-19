package domain

type PullRequestStore interface {
	// FindMatchingPullRequest returns the PullRequest that has the same branch as GitRepository.CommitBranch.
	// If not found, it returns nil without error.
	FindMatchingPullRequest(repository *GitRepository) (*PullRequest, error)

	// EnsurePullRequest creates or updates the given pull request in the repository.
	//  * This operation does not alter any properties of existing labels.
	//  * Existing labels are left untouched, but any extraneous labels are removed.
	//  * Title and Body are updated.
	//  * Existing Commit and Base branches are left untouched.
	//
	// The first error encountered aborts the operation.
	EnsurePullRequest(repository *GitRepository, pr *PullRequest) error
}
