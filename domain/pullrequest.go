package domain

import (
	"fmt"
)

type PullRequest struct {
	number       *PullRequestNumber
	title        string
	body         string
	CommitBranch string
	BaseBranch   string

	labels LabelSet
}

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

func (pr *PullRequest) GetLabels() LabelSet {
	return pr.labels
}

func (pr *PullRequest) SetNumber(nr *PullRequestNumber) error {
	if nr != nil && *nr <= 0 {
		return fmt.Errorf("%w: PR number cannot be lower than 1", ErrInvalidArgument)
	}
	pr.number = nr
	return nil
}

func (pr *PullRequest) GetNumber() *PullRequestNumber {
	return pr.number
}

func (pr *PullRequest) GetTitle() string {
	return pr.title
}

func (pr *PullRequest) GetBody() string {
	return pr.body
}

func (pr *PullRequest) ChangeDescription(title, body string) error {
	if err := pr.validateTitle(title); hasFailed(err) {
		return err
	}
	pr.title = title
	pr.body = body
	return nil
}

func (pr *PullRequest) AttachLabels(labels LabelSet) error {
	if err := pr.validateLabels(labels); hasFailed(err) {
		return err
	}
	pr.labels = labels
	return nil
}
