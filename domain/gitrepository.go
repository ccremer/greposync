package domain

// GitRepository is the heart of the domain.
//
// The model itself doesn't feature common actions like Commit.
// It was decided against adding those rich functionalities since that would mean implementing a replayable history of actions to keep in memory.
// This was considered too complicated, thus these actions are to be implemented in Stores.
type GitRepository struct {
	// RootDir is the full path to the Git root directory in the local filesystem.
	RootDir Path
	// URL is the remote URL of origin.
	URL *GitURL
	// PullRequest is the associated PullRequest for this repository in the remote Git hosting service.
	PullRequest *PullRequest
	// Labels contains the LabelSet that is present in the remote Git hosting service.
	Labels LabelSet

	// CommitBranch in the branch name of the current branch the working tree is in.
	CommitBranch string
	// DefaultBranch is the branch name of the remote default branch (usually `master` or `main`).
	DefaultBranch string
}

// NewGitRepository creates a new instance.
func NewGitRepository(u *GitURL, root Path) *GitRepository {
	return &GitRepository{
		URL:     u,
		RootDir: root,
	}
}

func (r *GitRepository) validateLabels(labels LabelSet) error {
	return firstOf(labels.CheckForEmptyLabelNames(), labels.CheckForDuplicates())
}

// SetLabels validates and sets the new LabelSet.
// Returns nil if there are no empty Label names or duplicates.
func (r *GitRepository) SetLabels(labels LabelSet) error {
	if err := r.validateLabels(labels); hasFailed(err) {
		return err
	}
	r.Labels = labels
	return nil
}

// AsValues returns the metadata as Values for rendering.
func (r GitRepository) AsValues() Values {
	return Values{
		"FullName":      r.URL.GetFullName(),
		"Name":          r.URL.GetRepositoryName(),
		"Namespace":     r.URL.GetNamespace(),
		"CommitBranch":  r.CommitBranch,
		"DefaultBranch": r.DefaultBranch,
		"RootDir":       r.RootDir,
	}
}
