package github

import (
	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v39/github"
)

type PrConverter struct{}

// ConvertToEntity converts the given object to another.
func (c PrConverter) ConvertToEntity(pr *github.PullRequest) *domain.PullRequest {
	if pr == nil {
		return nil
	}

	set := LabelSetConverter{}.ConvertToEntity(pr.Labels)
	// TODO: At least log a warning.
	// We don't expect invalid colors if coming from a repository, but that's just an assumption

	entity, _ := domain.NewPullRequest(domain.NewPullRequestNumber(pr.Number), *pr.Title, *pr.Body, *pr.Head.Ref, *pr.Base.Ref, set)

	return entity
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (c PrConverter) ConvertFromEntity(entity *domain.PullRequest) *github.PullRequest {
	if entity == nil {
		return nil
	}
	pr := &github.PullRequest{
		Number: entity.GetNumber().Int(),
		Title:  github.String(entity.GetTitle()),
		Body:   github.String(entity.GetBody()),
		Labels: LabelSetConverter{}.ConvertFromEntity(entity.GetLabels()),
	}
	return pr
}
