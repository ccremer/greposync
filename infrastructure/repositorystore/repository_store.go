package repositorystore

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	giturls "github.com/whilp/git-urls"
)

type RepositoryStore struct {
	StoreConfig
	k               *koanf.Koanf
	instrumentation *RepositoryStoreInstrumentation
}

// ManagedGitRepo is the representation of the managed git repos in the config file.
type ManagedGitRepo struct {
	Name string
}

type StoreConfig struct {
	ParentDir        string
	DefaultNamespace string
	CommitBranch     string

	IncludeFilter string
	ExcludeFilter string
	// ManagedReposFileName is the base file name where managed git repositories config is searched.
	ManagedReposFileName string
}

func NewRepositoryStore(instrumentation *RepositoryStoreInstrumentation) *RepositoryStore {
	return &RepositoryStore{
		k:               koanf.New("."),
		instrumentation: instrumentation,
		StoreConfig: StoreConfig{
			ManagedReposFileName: "managed_repos.yml",
		},
	}
}

func (s *RepositoryStore) deleteFileHandler(file *os.File) {
	_ = file.Close()
	_ = os.Remove(file.Name())
}

func (s *RepositoryStore) FetchGitRepositories() ([]*domain.GitRepository, error) {
	s.instrumentation.loadRepositoryConfigFile(s.ManagedReposFileName)
	if err := s.k.Load(file.Provider(s.ManagedReposFileName), yaml.Parser()); err != nil {
		return nil, err
	}
	var list []*domain.GitRepository
	var m []ManagedGitRepo
	if err := s.k.Unmarshal("repositories", &m); err != nil {
		return nil, err
	}
	gitBase := "git@github.com:"

	includeRegex, excludeRegex, err := compileRegex(s.IncludeFilter, s.ExcludeFilter)
	if err != nil {
		return nil, err
	}

	for _, repo := range m {
		u, err := parseUrl(repo, gitBase, s.DefaultNamespace)
		if err != nil {
			return list, err
		}

		gitUrl := domain.FromURL(u)
		if skipRepository(gitUrl.String(), includeRegex, excludeRegex) {
			s.instrumentation.skipRepository(gitUrl)
			continue
		}

		root := domain.NewFilePath(s.toLocalFilePath(gitUrl.AsURL()))
		gitRepository := domain.NewGitRepository(gitUrl, root)
		gitRepository.CommitBranch = s.CommitBranch
		if root.DirExists() {
			defaultBranch, err := GetDefaultBranch(gitRepository)
			if err != nil && !strings.Contains(err.Error(), "no default branch determined") {
				return list, err
			}
			gitRepository.DefaultBranch = defaultBranch
		}
		list = append(list, gitRepository)
	}
	return list, nil
}

func (s *RepositoryStore) toLocalFilePath(u *url.URL) string {
	p := strings.ReplaceAll(u.Path, "/", string(filepath.Separator))
	return filepath.Clean(filepath.Join(s.ParentDir, strings.ReplaceAll(u.Hostname(), ":", "-"), p))
}

func parseUrl(m ManagedGitRepo, gitBase, defaultNs string) (*url.URL, error) {
	if strings.Contains(m.Name, "/") {
		u, err := giturls.Parse(fmt.Sprintf("%s/%s", gitBase, m.Name))
		return u, err
	}
	u, err := giturls.Parse(fmt.Sprintf("%s/%s/%s", gitBase, defaultNs, m.Name))
	return u, err
}

func skipRepository(s string, includeRegex *regexp.Regexp, excludeRegex *regexp.Regexp) bool {
	return excludeRegex.MatchString(s) || !includeRegex.MatchString(s)
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
