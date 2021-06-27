package github

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

func CreatePR(r *git.Repository) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	err := createPR(context.Background(), client)
	CheckIfError(err)
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func createPR(ctx context.Context, client *github.Client) (err error) {
	prSubject := "Update from gsync"
	prRepo := "git-repo-sync"
	prRepoOwner := "ccremer"

	commitBranch := "my-branch"

	prBranch := "master"
	prDescription := "long text for PR desc"

	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &commitBranch,
		Base:                &prBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, prRepoOwner, prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}


func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
