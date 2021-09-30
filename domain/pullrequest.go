package domain

import (
	"fmt"
)

// PullRequest is a model that represents a pull request in a remote Git hosting service.
type PullRequest struct {
	number *PullRequestNumber
	title  string
	body   string

	// CommitBranch is the branch name of the current working tree.
	CommitBranch string
	// BaseBranch is the branch name into which CommitBranch should be merged into.
	BaseBranch string

	labels LabelSet
}

// NewPullRequest returns a new instance.
// An error is returned if the given properties do not satisfy constraints.
func NewPullRequest(
	number *PullRequestNumber, title, body, commitBranch, baseBranch string,
	labels LabelSet,
) (*PullRequest, error) {
	pr := &PullRequest{
		CommitBranch: commitBranch,
		BaseBranch:   baseBranch,
	}
	if err := firstOf(pr.SetNumber(number), pr.ChangeDescription(title, body), pr.AttachLabels(labels)); hasFailed(err) {
		return &PullRequest{}, err
	}
	return pr, nil
}

func (pr *PullRequest) validateLabels(labels LabelSet) error {
	return firstOf(labels.CheckForEmptyLabelNames(), labels.CheckForDuplicates())
}

func (pr *PullRequest) validateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidArgument)
	}
	return nil
}

// GetLabels returns the LabelSet of this PR.
func (pr *PullRequest) GetLabels() LabelSet {
	return pr.labels
}

// SetNumber sets the pull request number.
func (pr *PullRequest) SetNumber(nr *PullRequestNumber) error {
	if nr != nil && *nr <= 0 {
		return fmt.Errorf("%w: PR number cannot be lower than 1", ErrInvalidArgument)
	}
	pr.number = nr
	return nil
}

// GetNumber returns the pull request number.
// It returns nil if this PullRequest does not yet exist in remote.
func (pr *PullRequest) GetNumber() *PullRequestNumber {
	return pr.number
}

// GetTitle returns the PR title.
func (pr *PullRequest) GetTitle() string {
	return pr.title
}

// GetBody returns the PR description.
func (pr *PullRequest) GetBody() string {
	return pr.body
}

// ChangeDescription changes the title and body of this PR.
// An error is returned if the title is empty.
func (pr *PullRequest) ChangeDescription(title, body string) error {
	if err := pr.validateTitle(title); hasFailed(err) {
		return err
	}
	pr.title = title
	pr.body = body
	return nil
}

// AttachLabels sets the LabelSet of this PR.
// There cannot be duplicates or labels with no name.
func (pr *PullRequest) AttachLabels(labels LabelSet) error {
	if err := pr.validateLabels(labels); hasFailed(err) {
		return err
	}
	pr.labels = labels
	return nil
}
