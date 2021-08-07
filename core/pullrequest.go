package core

// PullRequest (also called Merge Request) is a request for change managed in the remote Git repository.
// The implementation may contain additional provider-specific properties.
//counterfeiter:generate . PullRequest
type PullRequest interface {
	GetTitle() string
	SetTitle(subject string)

	GetCommitBranch() string
	SetCommitBranch(string) string

	GetTargetBranch() string
	SetTargetBranch(string)

	GetBody() string
	SetBody(string)

	GetLabels() []Label
	SetLabels([]Label)
}
