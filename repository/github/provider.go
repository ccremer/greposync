package github

import (
	"context"
	"sync/atomic"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type (
	// Config configures the GitHub provider with all necessary information.
	Config struct {
		Token     string
		Repo      string
		RepoOwner string
	}
	// Provider contains the methods and data to interact with the GitHub API.
	Provider struct {
		cfg               *Config
		client            *github.Client
		ctx               context.Context
		log               printer.Printer
		remainingApiCalls int64
	}
)

func (p *Provider) convertLabelsToRepoLabels(labels []*github.Label) []cfg.RepositoryLabel {
	var list []cfg.RepositoryLabel
	for _, label := range labels {
		list = append(list, cfg.RepositoryLabel{
			Name:        label.GetName(),
			Color:       label.GetColor(),
			Description: label.GetDescription(),
		})
	}
	return list
}

// NewProvider returns a new GitHub provider instance.
func NewProvider(config *Config) *Provider {
	ctx, client := createClient(config.Token)
	provider := &Provider{
		cfg:    config,
		client: client,
		ctx:    ctx,
		log:    printer.New().SetName(config.Repo).SetLevel(printer.DefaultLevel),
	}
	return provider
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

func (p *Provider) setRemainingApiCalls(resp *github.Response) {
	if resp != nil {
		atomic.StoreInt64(&p.remainingApiCalls, int64(resp.Rate.Remaining))
	}
}

// GetRemainingApiCalls returns the amount of remaining rate-limited API request based on the last API request made.
func (p *Provider) GetRemainingApiCalls() int {
	return int(atomic.LoadInt64(&p.remainingApiCalls))
}
