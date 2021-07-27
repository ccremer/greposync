package github

import (
	"context"
	"os"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type (
	// Facade contains the methods and data to interact with the GitHub API.
	Facade struct {
		client *github.Client
		ctx    context.Context
		log    printer.Printer
	}
)

const GitHubProviderKey core.GitHostingProvider = "github"

// NewFacade returns a new GitHub provider instance.
func NewFacade() *Facade {
	provider := &Facade{
		log: printer.New(),
	}
	return provider
}

func (p *Facade) Initialize() error {
	ctx, client := createClient(os.Getenv("GITHUB_TOKEN"))
	p.client = client
	p.ctx = ctx
	return nil
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
