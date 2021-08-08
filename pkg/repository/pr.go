package repository

import "github.com/ccremer/greposync/core"

type PullRequestProperties struct {
	Title        string
	CommitBranch string
	TargetBranch string
	Body         string
}

// FetchPullRequest implements core.GitRepository.
func (s *Repository) FetchPullRequest() (core.PullRequest, error) {
	return s.remote.FindPullRequest(core.FromURL(s.GitConfig.Url), PullRequestProperties{
		CommitBranch: s.GitConfig.CommitBranch,
	})
}

// NewPullRequest implements core.GitRepository.
func (s *Repository) NewPullRequest() core.PullRequest {
	return s.remote.NewPullRequest(core.FromURL(s.GitConfig.Url), PullRequestProperties{
		CommitBranch: s.GitConfig.CommitBranch,
		TargetBranch: s.PrConfig.TargetBranch,
		Title:        s.PrConfig.Subject,
		Body:         s.PrConfig.BodyTemplate,
	})
}

// EnsurePullRequest implements core.GitRepository.
func (s *Repository) EnsurePullRequest(pr core.PullRequest) error {
	return s.remote.EnsurePullRequest(core.FromURL(s.GitConfig.Url), pr)
}
