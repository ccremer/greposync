package repository

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/go-git/go-git/v5"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	giturls "github.com/whilp/git-urls"
)

var k = koanf.New(".")

type (
	Service struct {
		r      *git.Repository
		p      printer.Printer
		Config Config
	}
	Config struct {
		GitDir     string
		GitUrl     string
		SkipReset  bool
		SkipCommit bool
		SkipPush   bool
		ForcePush  bool
		CreatePR   bool
		Amend      bool
		CommitMessage string
	}
	ManagedGitRepo struct {
		Name string
	}
)

func NewServicesFromFile(cfg *cfg.Configuration) []*Service {
	err := k.Load(file.Provider("managed_repos.yml"), yaml.Parser())
	printer.CheckIfError(err)

	var s []*Service
	var m []ManagedGitRepo
	err = k.Unmarshal("repositories", &m)
	printer.CheckIfError(err)
	gitBase := "git@github.com:"
	for _, repo := range m {
		u := parseUrl(repo, gitBase, cfg.Namespace)
		repoName := path.Base(u.Path)
		s = append(s, &Service{
			p: printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetLevel(printer.LevelDebug).SetName(repoName),
			Config: Config{
				GitDir:     path.Clean(path.Join(cfg.ProjectRoot, strings.ReplaceAll(u.Hostname(), ":", "-"), u.Path)),
				GitUrl:     u.String(),
				ForcePush:  true,
				SkipReset:  cfg.SkipReset,
				SkipPush:   cfg.SkipPush,
				SkipCommit: cfg.SkipCommit,
				CreatePR:   cfg.PullRequest.Create,
				Amend:      cfg.Amend,
				CommitMessage: cfg.Message,
			},
		})
	}
	return s
}

func parseUrl(m ManagedGitRepo, gitBase, defaultNs string) *url.URL {

	if strings.Contains(m.Name, "/") {
		u, err := giturls.Parse(fmt.Sprintf("%s/%s", gitBase, m.Name))
		printer.CheckIfError(err)
		return u
	}
	u, err := giturls.Parse(fmt.Sprintf("%s/%s/%s", gitBase, defaultNs, m.Name))
	printer.CheckIfError(err)
	return u
}
