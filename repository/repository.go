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

type (
	Service struct {
		r      *git.Repository
		p      printer.Printer
		Config cfg.GitConfig
	}
	ManagedGitRepo struct {
		Name string
	}
)

var (
	k                    = koanf.New(".")
	ManagedReposFileName = "managed_repos.yml"
)

func NewServicesFromFile(config *cfg.Configuration) []*Service {
	err := k.Load(file.Provider(ManagedReposFileName), yaml.Parser())
	printer.CheckIfError(err)

	var s []*Service
	var m []ManagedGitRepo
	err = k.Unmarshal("repositories", &m)
	printer.CheckIfError(err)
	gitBase := "git@github.com:"
	for _, repo := range m {
		u := parseUrl(repo, gitBase, config.Git.Namespace)
		repoName := path.Base(u.Path)
		s = append(s, &Service{
			p: printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetLevel(printer.LevelDebug).SetName(repoName),
			Config: cfg.GitConfig{
				Dir:           path.Clean(path.Join(config.ProjectRoot, strings.ReplaceAll(u.Hostname(), ":", "-"), u.Path)),
				Url:           u.String(),
				ForcePush:     true,
				SkipReset:     config.Git.SkipReset,
				SkipPush:      config.Git.SkipPush,
				SkipCommit:    config.Git.SkipCommit,
				Amend:         config.Git.Amend,
				CommitMessage: config.Git.CommitMessage,
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
