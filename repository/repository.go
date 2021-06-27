package repository

import (
	"fmt"
	"net/url"
	"path"
	"strings"

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
	}
	ManagedGitRepo struct {
		Name string
	}
)

func NewServicesFromFile(managedRepoPath string, repoRootDir string, defaultNs string) []*Service {
	err := k.Load(file.Provider(managedRepoPath), yaml.Parser())
	printer.CheckIfError(err)

	var s []*Service
	var m []ManagedGitRepo
	err = k.Unmarshal("repositories", &m)
	printer.CheckIfError(err)
	gitBase := "git@github.com:"
	for _, repo := range m {
		u := parseUrl(repo, gitBase, defaultNs)
		s = append(s, &Service{
			p: printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetLevel(printer.LevelDebug),
			Config: Config{
				GitDir:     path.Clean(path.Join(repoRootDir, strings.ReplaceAll(u.Hostname(), ":", "-"), u.Path)),
				GitUrl:     u.String(),
				ForcePush:  true,
				SkipReset:  true,
				SkipPush:   true,
				SkipCommit: true,
				CreatePR:   false,
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
