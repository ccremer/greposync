package github

import (
	"context"
	"fmt"
	"log"

	"github.com/ccremer/git-repo-sync/printer"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

type (
	Config struct {
		Token        string
		Subject      string
		Repo         string
		RepoOwner    string
		CommitBranch string
		TargetBranch string
		Body         string
		Labels       []string
	}
)

func CreatePR(c Config) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	err := createPR(context.Background(), client, c)
	printer.CheckIfError(err)
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR(ctx context.Context, client *github.Client, c Config) (err error) {
	newPR := &github.NewPullRequest{
		Title:               &c.Subject,
		Head:                &c.CommitBranch,
		Base:                &c.TargetBranch,
		Body:                &c.Body,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, c.RepoOwner, c.Repo, newPR)
	if err != nil {
		return err
	}

	if len(c.Labels) > 0 {
		addLabels(ctx, client, c, *pr.Number)
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func addLabels(ctx context.Context, client *github.Client, c Config, issueNumber int) {
	_, _, err := client.Issues.AddLabelsToIssue(ctx, c.RepoOwner, c.Repo, issueNumber, c.Labels)
	if err != nil {
		log.Println("could not add label, ignoring error: " + err.Error())
	}
}
