package github

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/printer"
	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

type (
	// GhRemote contains the methods and data to interact with the GitHub API.
	GhRemote struct {
		client *github.Client
		ctx    context.Context
		log    printer.Printer
		m      *sync.Mutex
	}
)

// GitHubProviderKey is the identifier for the GitHub core.GitHostingProvider.
const GitHubProviderKey core.GitHostingProvider = "github"

// NewRemote returns a new GitHub provider instance.
func NewRemote() *GhRemote {
	provider := &GhRemote{
		log: printer.New(),
		m:   &sync.Mutex{},
	}
	return provider
}

// Initialize implements repository.Remote.
func (r *GhRemote) Initialize() error {
	ctx, client := createClient(os.Getenv("GITHUB_TOKEN"))
	r.client = client
	r.ctx = ctx
	return nil
}

// FindLabels implements repository.Remote.
func (r *GhRemote) FindLabels(url *core.GitURL, labels []*cfg.RepositoryLabel) ([]core.Label, error) {
	ghLabels, err := r.fetchAllLabels(url)
	if err != nil {
		return []core.Label{}, err
	}
	var impls = make([]*LabelImpl, len(labels))
	for i, configLabel := range labels {
		impl := &LabelImpl{
			Name:        configLabel.Name,
			Description: configLabel.Description,
			Color:       configLabel.Color,
			Inactive:    configLabel.Delete,
			remote:      r,
			repo:        url,
		}
		ghLabel := r.findMatchingGhLabel(ghLabels, impl)
		impl.ghLabel = ghLabel
		impls[i] = impl
	}
	return LabelConverter{}.ConvertToEntity(impls), nil
}

// FindPullRequest implements repository.Remote.
func (r *GhRemote) FindPullRequest(url *core.GitURL, config repository.PullRequestProperties) (core.PullRequest, error) {
	pr, err := r.findExistingPr(url.GetNamespace(), url.GetRepositoryName(), config.CommitBranch)
	if err != nil {
		return nil, err
	}
	return &GhPullRequest{
		PullRequestProperties: config,
		Owner:                 url.GetNamespace(),
		Repository:            url.GetRepositoryName(),
		client:                r.client,
		ghPullRequest:         pr,
		log:                   printer.New().SetName(url.GetRepositoryName()),
	}, nil
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

func (r *GhRemote) NewPullRequest(url *core.GitURL, config repository.PullRequestProperties) core.PullRequest {
	return &GhPullRequest{
		PullRequestProperties: config,
		Owner:                 url.GetNamespace(),
		Repository:            url.GetRepositoryName(),
		client:                r.client,
	}
}

func (r *GhRemote) EnsurePullRequest(_ *core.GitURL, p core.PullRequest) error {
	pr := p.(*GhPullRequest)
	if pr.ghPullRequest == nil {
		return pr.create()
	}
	return pr.update()
}

func createClient(token string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client
}
