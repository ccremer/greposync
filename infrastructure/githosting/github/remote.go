package github

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type (
	// GhRemote contains the methods and data to interact with the GitHub API.
	GhRemote struct {
		client          *github.Client
		ctx             context.Context
		m               *sync.Mutex
		prCache         map[int]*github.PullRequest
		labelCache      map[*domain.GitURL][]*github.Label
		instrumentation *GitHubInstrumentation
	}
)

// ProviderKey is the identifier for the GitHub githosting.RemoteProvider.
const ProviderKey githosting.RemoteProvider = "github"

// NewRemote returns a new GitHub provider instance.
func NewRemote(instrumentation *GitHubInstrumentation) *GhRemote {
	ctx := context.Background()
	provider := &GhRemote{
		m:               &sync.Mutex{},
		ctx:             ctx,
		client:          createClient(os.Getenv("GITHUB_TOKEN"), ctx),
		prCache:         map[int]*github.PullRequest{},
		labelCache:      map[*domain.GitURL][]*github.Label{},
		instrumentation: instrumentation,
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

// delayedUnlock sleeps one second for abuse rate limit best-practice and releases the lock.
//
// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
// "If you're making a large number of POST, PATCH, PUT, or DELETE requests for a single user or client ID, wait at least one second between each request."
func (r *GhRemote) delayedUnlock() {
	time.Sleep(1 * time.Second)
	r.m.Unlock()
}
