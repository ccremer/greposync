package repository

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	giturls "github.com/whilp/git-urls"
)

type (
	// Service represents a git repository that comes with utility methods
	Service struct {
		p      printer.Printer
		Config *cfg.GitConfig
	}
	// ManagedGitRepo is the representation of the managed git repos in the config file.
	ManagedGitRepo struct {
		Name string
	}
)

var (
	k = koanf.New(".")
	// ManagedReposFileName is the base file name where managed git repositories config is searched.
	ManagedReposFileName = "managed_repos.yml"
)

// NewServicesFromFile parses a config file with managed git repositories and provides a Service for each.
func NewServicesFromFile(config *cfg.Configuration) []*Service {
	err := k.Load(file.Provider(ManagedReposFileName), yaml.Parser())
	printer.CheckIfError(err)

	var list []*Service
	var m []ManagedGitRepo
	err = k.Unmarshal("repositories", &m)
	printer.CheckIfError(err)
	gitBase := "git@github.com:"
	for _, repo := range m {
		u := parseUrl(repo, gitBase, config.Git.Namespace)
		repoName := path.Base(u.Path)
		s := &Service{
			p: printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetLevel(printer.DefaultLevel).SetName(repoName),
			Config: &cfg.GitConfig{
				Dir:           path.Clean(path.Join(config.ProjectRoot, strings.ReplaceAll(u.Hostname(), ":", "-"), u.Path)),
				Url:           u,
				ForcePush:     true,
				SkipReset:     config.Git.SkipReset,
				SkipPush:      config.Git.SkipPush,
				SkipCommit:    config.Git.SkipCommit,
				Amend:         config.Git.Amend,
				CommitMessage: config.Git.CommitMessage,
				CommitBranch:  config.Git.CommitBranch,
				Namespace:     config.Git.Namespace,
				Name:          repoName,
			},
		}
		list = append(list, s)
	}
	return list
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

// DirExists returns true if the given path exists and is a directory.
func (s *Service) DirExists(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
