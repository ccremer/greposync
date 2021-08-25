package github

import (
	"context"
	"os"
	"sync"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/githosting"
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
		m       *sync.Mutex
		prCache map[int]*github.PullRequest
	}
)

// ProviderKey is the identifier for the GitHub githosting.RemoteProvider.
const ProviderKey githosting.RemoteProvider = "github"

// NewRemote returns a new GitHub provider instance.
func NewRemote() *GhRemote {
	ctx := context.Background()
	provider := &GhRemote{
		log:     printer.New(),
		m:       &sync.Mutex{},
		ctx:     ctx,
		client:  createClient(os.Getenv("GITHUB_TOKEN"), ctx),
		prCache: map[int]*github.PullRequest{},
	}
	return provider
}

func createClient(token string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func (r *GhRemote) HasSupportFor(url *domain.GitURL) bool {
	return url.Host == "github.com"
}
