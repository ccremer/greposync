package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v38/github"
)

type (
	// PrConfig is a config required for managing PRs.
	PrConfig struct {
		Subject      string
		CommitBranch string
		TargetBranch string
		Body         string
		Labels       []string
	}
)

// CreateOrUpdatePr creates the PR if it doesn't exist, or updates an existing one if the branch matches.
// A PR is considered out-of-date if the subject or body don't match with current configuration.
// Labels are left unmodified.
func (p *Provider) CreateOrUpdatePr(c *PrConfig) error {
	if pr, err := p.findExistingPr(c); err != nil {
		return err
	} else if pr == nil {
		return p.createPR(c)
	} else {
		if *pr.Body != c.Body || *pr.Title != c.Subject {
			p.log.InfoF("Updating PR#%d", *pr.Number)
			return p.updatePr(c, pr)
		} else {
			p.log.InfoF("PR#%d is up-to-date.", *pr.Number)
			return nil
		}
	}
}

func (p *Provider) findExistingPr(c *PrConfig) (*github.PullRequest, error) {
	p.log.DebugF("Searching existing PRs with same branch %s...", c.CommitBranch)
	list, resp, err := p.client.PullRequests.List(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, &github.PullRequestListOptions{
		Head: fmt.Sprintf("%s:%s", p.cfg.RepoOwner, c.CommitBranch),
	})
	if err != nil {
		return nil, err
	}
	p.setRemainingApiCalls(resp)
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (p *Provider) createPR(c *PrConfig) (err error) {
	p.log.DebugF("Creating new PR")
	newPR := &github.NewPullRequest{
		Title:               &c.Subject,
		Head:                &c.CommitBranch,
		Base:                &c.TargetBranch,
		Body:                &c.Body,
		MaintainerCanModify: github.Bool(true),
	}

	pr, resp, err := p.client.PullRequests.Create(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, newPR)
	if err != nil {
		if strings.Contains(err.Error(), "No commits between") {
			p.log.InfoF("No pull request created as there are no commits between '%s' and '%s'", c.TargetBranch, c.CommitBranch)
			return nil
		}
		return err
	}
	p.setRemainingApiCalls(resp)

	if len(c.Labels) > 0 {
		p.addLabels(c, *pr.Number)
	}

	p.log.InfoF("PR created: %s", pr.GetHTMLURL())
	return nil
}

func (p *Provider) addLabels(c *PrConfig, issueNumber int) {
	_, resp, err := p.client.Issues.AddLabelsToIssue(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, issueNumber, c.Labels)
	if err != nil {
		p.log.WarnF("could not add label, ignoring error: " + err.Error())
	}
	p.setRemainingApiCalls(resp)
}

func (p *Provider) updatePr(c *PrConfig, pr *github.PullRequest) error {
	pr.Body = &c.Body
	pr.Title = &c.Subject
	_, resp, err := p.client.PullRequests.Edit(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, *pr.Number, pr)
	p.setRemainingApiCalls(resp)
	return err
}
