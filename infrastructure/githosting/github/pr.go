package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v39/github"
)

func (r *GhRemote) FindPullRequest(repository *domain.GitRepository) (*domain.PullRequest, error) {
	pr, err := r.findExistingPr(repository)
	if err != nil {
		return nil, err
	}
	if pr != nil {
		r.prCache[*pr.Number] = pr
	}
	converted := PrConverter{}.ConvertToEntity(pr)
	return converted, nil
}

func (r *GhRemote) findExistingPr(repository *domain.GitRepository) (*github.PullRequest, error) {
	list, _, err := r.client.PullRequests.List(context.Background(), repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), &github.PullRequestListOptions{
		Head: fmt.Sprintf("%s:%s", repository.URL.GetNamespace(), repository.CommitBranch),
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], r.instrumentation.prFound(repository, list[0])
	}
	return nil, r.instrumentation.noPrFound(repository)
}

func (r *GhRemote) EnsurePullRequest(repository *domain.GitRepository, pr *domain.PullRequest) error {
	converted := PrConverter{}.ConvertFromEntity(pr)
	cached, exists := r.prCache[converted.GetNumber()]
	if !exists {
		return r.createNewPr(repository, pr)
	}
	cached.Title = converted.Title
	cached.Body = converted.Body
	cached.Labels = converted.Labels
	return r.updateExistingPr(repository, cached, pr)
}

func (r *GhRemote) updateExistingPr(repository *domain.GitRepository, cached *github.PullRequest, pr *domain.PullRequest) error {
	if r.canSkipDescriptionUpdate(cached, pr) && r.canSkipLabelUpdate(cached, pr) {
		return r.instrumentation.prIsUpToDate(repository, cached)
	}
	err := r.updatePrDescription(repository, cached, pr)
	if err != nil {
		return err
	}
	return r.updatePrLabels(repository, cached, pr)
}

func (r *GhRemote) updatePrDescription(repository *domain.GitRepository, cached *github.PullRequest, pr *domain.PullRequest) error {
	if r.canSkipDescriptionUpdate(cached, pr) {
		return nil
	}
	ghPr, _, err := r.client.PullRequests.Edit(context.Background(), repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), *cached.Number, cached)
	return r.instrumentation.prUpdated(repository, ghPr, err)
}

func (r *GhRemote) updatePrLabels(repository *domain.GitRepository, cached *github.PullRequest, pr *domain.PullRequest) error {
	if r.canSkipLabelUpdate(cached, pr) {
		return nil
	}
	err := r.setLabelsToPr(repository.URL, *cached.Number, pr.GetLabels())
	return r.instrumentation.prLabelsUpdated(repository, pr, err)
}

func (r *GhRemote) canSkipDescriptionUpdate(cached *github.PullRequest, pr *domain.PullRequest) bool {
	sameTitle := *cached.Title == pr.GetTitle()
	sameBody := *cached.Body == pr.GetBody()
	return sameTitle && sameBody
}

func (r *GhRemote) canSkipLabelUpdate(cached *github.PullRequest, pr *domain.PullRequest) bool {
	converted := LabelSetConverter{}.ConvertToEntity(cached.Labels)
	diff := pr.GetLabels().Without(converted)
	return len(diff) == 0
}

// createNewPr makes a new pull request in GitHub.
// Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (r *GhRemote) createNewPr(repository *domain.GitRepository, pr *domain.PullRequest) error {
	newPR := &github.NewPullRequest{
		Title:               github.String(pr.GetTitle()),
		Head:                &pr.CommitBranch,
		Base:                &pr.BaseBranch,
		Body:                github.String(pr.GetBody()),
		MaintainerCanModify: github.Bool(true),
	}

	ghPr, _, err := r.client.PullRequests.Create(context.Background(), repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), newPR)
	if err != nil {
		if strings.Contains(err.Error(), "No commits between") {
			return r.instrumentation.prNotCreatedBecauseNoCommits(repository, pr)
		}
		return err
	}

	if len(pr.GetLabels()) > 0 {
		err := r.setLabelsToPr(repository.URL, *ghPr.Number, pr.GetLabels())
		if err != nil {
			return err
		}
	}

	r.instrumentation.prCreated(repository, ghPr.GetHTMLURL())
	return nil
}

func (r *GhRemote) setLabelsToPr(url *domain.GitURL, issueNumber int, set domain.LabelSet) error {
	var labelArr = make([]string, len(set))
	for i := range set {
		labelArr[i] = set[i].Name
	}
	_, _, err := r.client.Issues.ReplaceLabelsForIssue(context.Background(), url.GetNamespace(), url.GetRepositoryName(), issueNumber, labelArr)
	return err
}
