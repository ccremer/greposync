package repository

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
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
func NewServicesFromFile(config *cfg.Configuration) ([]*Service, error) {
	if err := k.Load(file.Provider(ManagedReposFileName), yaml.Parser()); err != nil {
		return nil, err
	}
	var list []*Service
	var m []ManagedGitRepo
	if err := k.Unmarshal("repositories", &m); err != nil {
		return nil, err
	}
	gitBase := "git@github.com:"

	includeRegex, excludeRegex, err := compileRegex(config.Project.Include, config.Project.Exclude)
	if err != nil {
		return nil, err
	}

	for _, repo := range m {
		u := parseUrl(repo, gitBase, config.Git.Namespace)
		repoName := path.Base(u.Path)
		log := printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetLevel(printer.DefaultLevel).SetName(repoName)
		if skipRepository(u, includeRegex, excludeRegex) {
			log.InfoF("Skipping '%s%s' due to filters", u.Hostname(), u.Path)
			continue
		}

		s := &Service{
			p: log,
			Config: &cfg.GitConfig{
				Dir:           path.Clean(path.Join(config.Project.RootDir, strings.ReplaceAll(u.Hostname(), ":", "-"), u.Path)),
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
	return list, nil
}

func skipRepository(u *url.URL, includeRegex *regexp.Regexp, excludeRegex *regexp.Regexp) bool {
	return matchRegex(excludeRegex, u) || !matchRegex(includeRegex, u)
}

func compileRegex(includeFilter, excludeFilter string) (includeRegex, excludeRegex *regexp.Regexp, err error) {
	if includeFilter == "" {
		includeFilter = ".*"
	}
	if excludeFilter == "" {
		excludeFilter = "^$"
	}
	includeRegex, err = regexp.Compile(includeFilter)
	if err != nil {
		return nil, nil, err
	}
	excludeRegex, err = regexp.Compile(excludeFilter)
	if err != nil {
		return nil, nil, err
	}
	return includeRegex, excludeRegex, nil
}

func matchRegex(regex *regexp.Regexp, u *url.URL) bool {
	return regex.MatchString(u.String())
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
