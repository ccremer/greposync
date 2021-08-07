package github

import (
	"context"
	"strings"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/printer"
	"github.com/google/go-github/v37/github"
)

type GhPullRequest struct {
	repository.PullRequestProperties
	Owner         string
	Repository    string
	client        *github.Client
	ghPullRequest *github.PullRequest
	log           printer.Printer
	Labels        []core.Label
}

func (pr *GhPullRequest) GetTitle() string {
	return pr.Title
}

func (pr *GhPullRequest) SetTitle(title string) {
	pr.Title = title
}

func (pr *GhPullRequest) GetCommitBranch() string {
	return pr.CommitBranch
}

func (pr *GhPullRequest) SetCommitBranch(s string) string {
	panic("implement me")
}

func (pr *GhPullRequest) GetTargetBranch() string {
	panic("implement me")
}

func (pr *GhPullRequest) SetTargetBranch(s string) {
	panic("implement me")
}

func (pr *GhPullRequest) GetBody() string {
	return pr.Body
}

func (pr *GhPullRequest) SetBody(s string) {
	pr.Body = s
}

func (pr *GhPullRequest) GetLabels() []core.Label {
	return pr.Labels
}

func (pr *GhPullRequest) SetLabels(labels []core.Label) {
	pr.Labels = labels
}

// create makes a new pull request in GitHub. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (pr *GhPullRequest) create() error {
	pr.log.DebugF("Creating new PR")
	newPR := &github.NewPullRequest{
		Title:               &pr.Title,
		Head:                &pr.CommitBranch,
		Base:                &pr.TargetBranch,
		Body:                &pr.Body,
		MaintainerCanModify: github.Bool(true),
	}

	ghPr, _, err := pr.client.PullRequests.Create(context.Background(), pr.Owner, pr.Repository, newPR)
	if err != nil {
		if strings.Contains(err.Error(), "No commits between") {
			pr.log.InfoF("No pull request created as there are no commits between '%s' and '%s'", pr.TargetBranch, pr.CommitBranch)
			return nil
		}
		return err
	}
	pr.ghPullRequest = ghPr

	if len(pr.Labels) > 0 {
		pr.addLabelsToPr()
	}

	pr.log.InfoF("PR created: %s", ghPr.GetHTMLURL())
	return nil
}

func (pr *GhPullRequest) addLabelsToPr() {
	var labelArr = make([]string, len(pr.Labels))
	for i := range pr.Labels {
		labelArr[i] = pr.Labels[i].GetName()
	}
	_, _, err := pr.client.Issues.AddLabelsToIssue(context.Background(), pr.Owner, pr.Repository, *pr.ghPullRequest.Number, labelArr)
	if err != nil {
		pr.log.WarnF("could not add label, ignoring error: " + err.Error())
	}
}

func (pr *GhPullRequest) update() error {
	pr.ghPullRequest.Body = &pr.Body
	pr.ghPullRequest.Title = &pr.Title
	_, _, err := pr.client.PullRequests.Edit(context.Background(), pr.Owner, pr.Repository, *pr.ghPullRequest.Number, pr.ghPullRequest)
	return err
}
