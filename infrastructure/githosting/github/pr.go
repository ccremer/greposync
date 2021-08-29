package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v38/github"
)

func (r *GhRemote) FindPullRequest(url *domain.GitURL, _, commitBranch string) (*domain.PullRequest, error) {
	pr, err := r.findExistingPr(url.GetNamespace(), url.GetRepositoryName(), commitBranch)
	if err != nil {
		return nil, err
	}
	if pr != nil {
		r.prCache[*pr.Number] = pr
	}
	converted := PrConverter{}.ConvertToEntity(pr)
	return converted, nil
}

func (r *GhRemote) findExistingPr(owner, repo, commitBranch string) (*github.PullRequest, error) {
	r.log.DebugF("Searching existing PRs with same branch %s...", commitBranch)
	list, _, err := r.client.PullRequests.List(context.Background(), owner, repo, &github.PullRequestListOptions{
		Head: fmt.Sprintf("%s:%s", owner, commitBranch),
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

func (r *GhRemote) EnsurePullRequest(url *domain.GitURL, pr *domain.PullRequest) error {
	converted := PrConverter{}.ConvertFromEntity(pr)
	cached, exists := r.prCache[*converted.Number]
	if !exists {
		return r.createNewPr(url, pr)
	}
	cached.Title = converted.Title
	cached.Body = converted.Body
	cached.Labels = converted.Labels
	converted = cached
	return r.updateExistingPr(url, converted, pr)
}

func (r *GhRemote) updateExistingPr(url *domain.GitURL, cached *github.PullRequest, pr *domain.PullRequest) error {
	err := r.updatePrDescription(url, cached, pr)
	if err != nil {
		return err
	}
	return r.updatePrLabels(url, cached, pr)
}

func (r *GhRemote) updatePrDescription(url *domain.GitURL, cached *github.PullRequest, pr *domain.PullRequest) error {
	if r.canSkipDescriptionUpdate(cached, pr) {
		return nil
	}
	_, _, err := r.client.PullRequests.Edit(context.Background(), url.GetNamespace(), url.GetRepositoryName(), *cached.Number, cached)
	return err
}

func (r *GhRemote) updatePrLabels(url *domain.GitURL, cached *github.PullRequest, pr *domain.PullRequest) error {
	if r.canSkipLabelUpdate(cached, pr) {
		return nil
	}
	return r.setLabelsToPr(url, *cached.Number, pr.GetLabels())
}

func (r *GhRemote) canSkipDescriptionUpdate(cached *github.PullRequest, pr *domain.PullRequest) bool {
	sameTitle := *cached.Title == pr.Title
	sameBody := *cached.Body == pr.Body
	return sameTitle && sameBody
}

func (r *GhRemote) canSkipLabelUpdate(cached *github.PullRequest, pr *domain.PullRequest) bool {
	converted := LabelSetConverter{}.ConvertToEntity(cached.Labels)
	diff := pr.GetLabels().DifferenceOf(converted)
	return len(diff) == 0
}

// createNewPr makes a new pull request in GitHub.
// Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (r *GhRemote) createNewPr(url *domain.GitURL, pr *domain.PullRequest) error {
	newPR := &github.NewPullRequest{
		Title:               &pr.Title,
		Head:                &pr.CommitBranch,
		Base:                &pr.BaseBranch,
		Body:                &pr.Body,
		MaintainerCanModify: github.Bool(true),
	}

	ghPr, _, err := r.client.PullRequests.Create(context.Background(), url.GetNamespace(), url.GetRepositoryName(), newPR)
	if err != nil {
		if strings.Contains(err.Error(), "No commits between") {
			r.log.InfoF("No pull request created as there are no commits between '%s' and '%s'", pr.BaseBranch, pr.CommitBranch)
			return nil
		}
		return err
	}

	if len(pr.GetLabels()) > 0 {
		err := r.setLabelsToPr(url, *ghPr.Number, pr.GetLabels())
		if err != nil {
			return err
		}
	}

	r.log.InfoF("PR created: %s", ghPr.GetHTMLURL())
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
