package repository

import "github.com/ccremer/greposync/core"

type PullRequestProperties struct {
	Title        string
	CommitBranch string
	TargetBranch string
	Body         string
}

// FetchPullRequest implements core.GitRepository.
func (g *Repository) FetchPullRequest() (core.PullRequest, error) {
	return g.remote.FindPullRequest(core.FromURL(g.GitConfig.Url), PullRequestProperties{
		CommitBranch: g.GitConfig.CommitBranch,
	})
}

// NewPullRequest implements core.GitRepository.
func (g *Repository) NewPullRequest() core.PullRequest {
	return g.remote.NewPullRequest(core.FromURL(g.GitConfig.Url), PullRequestProperties{
		CommitBranch: g.GitConfig.CommitBranch,
		TargetBranch: g.PrConfig.TargetBranch,
		Title:        g.PrConfig.Subject,
		Body:         g.PrConfig.BodyTemplate,
	})
}

// EnsurePullRequest implements core.GitRepository.
func (g *Repository) EnsurePullRequest(pr core.PullRequest) error {
	return g.remote.EnsurePullRequest(core.FromURL(g.GitConfig.Url), pr)
}
