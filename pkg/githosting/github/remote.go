package github

import (
	"context"
	"os"
	"sync"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type (
	// Remote contains the methods and data to interact with the GitHub API.
	Remote struct {
		client *github.Client
		ctx    context.Context
		log    printer.Printer
		m      *sync.Mutex
	}
)

// GitHubProviderKey is the identifier for the GitHub core.GitHostingProvider.
const GitHubProviderKey core.GitHostingProvider = "github"

// NewRemote returns a new GitHub provider instance.
func NewRemote() *Remote {
	provider := &Remote{
		log: printer.New(),
		m:   &sync.Mutex{},
	}
	return provider
}

// Initialize implements core.GitHostingFacade.
func (p *Remote) Initialize() error {
	ctx, client := createClient(os.Getenv("GITHUB_TOKEN"))
	p.client = client
	p.ctx = ctx
	return nil
}

func (p *Remote) FindLabels(url *core.GitURL, labels []*cfg.RepositoryLabel) ([]core.Label, error) {
	ghLabels, err := p.fetchAllLabels(url)
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
			remote:      p,
			repo:        url,
		}
		ghLabel := p.findMatchingGhLabel(ghLabels, impl)
		impl.ghLabel = ghLabel
		impls[i] = impl
	}
	return LabelConverter{}.ConvertToEntity(impls), nil
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
