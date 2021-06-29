package github

import (
	"context"
	"fmt"

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
	PrProvider struct {
		cfg    *Config
		client *github.Client
		ctx    context.Context
		log    printer.Printer
	}
)

func (p *PrProvider) createClient() {
	p.ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.cfg.Token},
	)
	tc := oauth2.NewClient(p.ctx, ts)

	p.client = github.NewClient(tc)
}

func NewProvider(config *Config) *PrProvider {
	provider := &PrProvider{
		cfg: config,
	}
	provider.createClient()
	return provider
}

func (p *PrProvider) CreateOrUpdatePR() {

	if pr := p.findExistingPr(); pr == nil {
		err := p.createPR()
		p.log.CheckIfError(err)
	} else {
		if *pr.Body != p.cfg.Body || *pr.Title != p.cfg.Subject {
			p.log.InfoF("Updating PR#%s", pr.Number)
			err := p.updatePr(pr)
			p.log.CheckIfError(err)
		} else {
			p.log.InfoF("PR#%s is up-to-date.", pr.Number)
		}
	}

}

func (p *PrProvider) findExistingPr() *github.PullRequest {
	p.log.DebugF("Searching existing PRs with same branch %s...", p.cfg.CommitBranch)
	list, _, err := p.client.PullRequests.List(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, &github.PullRequestListOptions{
		Head: fmt.Sprintf("%s:%s", p.cfg.RepoOwner, p.cfg.CommitBranch),
	})
	p.log.CheckIfError(err)
	if len(list) > 0 {
		return list[0]
	}
	return nil
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (p *PrProvider) createPR() (err error) {
	p.log.DebugF("Creating new PR")
	newPR := &github.NewPullRequest{
		Title:               &p.cfg.Subject,
		Head:                &p.cfg.CommitBranch,
		Base:                &p.cfg.TargetBranch,
		Body:                &p.cfg.Body,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := p.client.PullRequests.Create(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, newPR)
	if err != nil {
		return err
	}

	if len(p.cfg.Labels) > 0 {
		p.addLabels(*pr.Number)
	}

	p.log.InfoF("PR created: %s", pr.GetHTMLURL())
	return nil
}

func (p *PrProvider) addLabels(issueNumber int) {
	_, _, err := p.client.Issues.AddLabelsToIssue(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, issueNumber, p.cfg.Labels)
	if err != nil {
		p.log.WarnF("could not add label, ignoring error: " + err.Error())
	}
}

func (p *PrProvider) updatePr(pr *github.PullRequest) error {
	pr.Body = &p.cfg.Body
	pr.Title = &p.cfg.Subject
	_, _, err := p.client.PullRequests.Edit(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, *pr.Number, pr)
	return err
}
