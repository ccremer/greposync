package domain

type GitRepository struct {
	RootDir      Path
	URL         *GitURL
	PullRequest *PullRequest
	Labels      LabelSet

	CommitBranch  string
	DefaultBranch string
}

func NewGitRepository(u *GitURL, root Path, labels LabelSet) (*GitRepository, error) {
	r := &GitRepository{
		URL:     u,
		RootDir: root,
	}
	if err := r.validateLabels(labels); hasFailed(err) {
		return r, err
	}
	r.Labels = labels
	return r, nil
}

func (r *GitRepository) validateLabels(labels LabelSet) error {
	return firstOf(labels.CheckForEmptyLabelNames(), labels.CheckForDuplicates())
}

func (r *GitRepository) SetLabels(labels LabelSet) error {
	if err := r.validateLabels(labels); hasFailed(err) {
		return err
	}
	r.Labels = labels
	return nil
}
