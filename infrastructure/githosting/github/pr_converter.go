package github

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v38/github"
)

type PrConverter struct{}

// ConvertToEntity converts the given object to another.
func (c PrConverter) ConvertToEntity(pr *github.PullRequest) *domain.PullRequest {
	if pr == nil {
		return nil
	}

	entity := &domain.PullRequest{
		Number:       c.toEntityNumber(pr.Number),
		Title:        *pr.Title,
		Body:         *pr.Body,
		CommitBranch: *pr.Head.Ref,
		BaseBranch:   *pr.Base.Ref,
	}

	set := LabelConverter{}.ConvertToEntity(pr.Labels)
	// TODO: Ignore for now
	_ = entity.AttachLabels(set)
	return entity
}

// ConvertFromEntity converts the given object to another.
// Returns a non-nil empty list if labels is empty or nil.
func (c PrConverter) ConvertFromEntity(entity *domain.PullRequest) *github.PullRequest {
	if entity == nil {
		return nil
	}
	pr := &github.PullRequest{
		Number: c.toGhNumber(entity.Number),
		Title:  &entity.Title,
		Body:   &entity.Body,
		Labels: LabelConverter{}.ConvertFromEntity(entity.GetLabels()),
	}
	return pr
}

func (PrConverter) toEntityNumber(number *int) string {
	if number == nil {
		return ""
	}
	return fmt.Sprintf("#%d", *number)
}

func (PrConverter) toGhNumber(number string) *int {
	if number == "" {
		return nil
	}
	raw := strings.TrimPrefix(number, "#")
	if nr, err := strconv.Atoi(raw); err != nil {
		return &nr
	}
	return nil
}
