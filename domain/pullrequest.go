package domain

import (
	"fmt"
)

type PullRequest struct {
	Number       string
	Title        string
	Body         string
	CommitBranch string
	BaseBranch   string

	labels LabelSet

	repo *GitRepository
}

func NewPullRequest(repo *GitRepository, labels LabelSet) (*PullRequest, error) {
	pr := &PullRequest{}
	if err := pr.validateLabels(labels); hasFailed(err) {
		return pr, err
	}
	pr.repo = repo
	pr.labels = labels
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

func (pr *PullRequest) GetRepository() *GitRepository {
	return pr.repo
}

func (pr *PullRequest) GetLabels() LabelSet {
	return pr.labels
}

func (pr *PullRequest) ChangeDescription(title, body string, store PullRequestStore) error {
	if err := pr.validateTitle(title); hasFailed(err) {
		return err
	}
	pr.Title = title
	pr.Body = body
	return nil
}

func (pr *PullRequest) AttachLabels(labels LabelSet) error {
	if err := pr.validateLabels(labels); hasFailed(err) {
		return err
	}
	pr.labels = labels
	return nil
}
