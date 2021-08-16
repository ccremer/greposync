package domain

type PullRequestStore interface {
	// SetLabelsInPullRequest ensures all labels exist in the given PullRequest.GetLabels.
	// Existing labels are untouched, but any extraneous labels are removed.
	// This operation does not alter any properties of existing labels.
	// Any error encountered aborts the operation immediately.
	SetLabelsInPullRequest(repository *GitRepository, pr *PullRequest) error


	//FindMatchingPullRequests(repository *GitRepository) ([]*PullRequest, error)
}
