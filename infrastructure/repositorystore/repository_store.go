package repositorystore

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	giturls "github.com/whilp/git-urls"
)

type RepositoryStore struct {
	StoreConfig
	log printer.Printer
	k   *koanf.Koanf
}

// ManagedGitRepo is the representation of the managed git repos in the config file.
type ManagedGitRepo struct {
	Name string
}

type StoreConfig struct {
	ParentDir        string
	DefaultNamespace string
}

var (
	// ManagedReposFileName is the base file name where managed git repositories config is searched.
	ManagedReposFileName = "managed_repos.yml"
)

func NewRepositoryStore() *RepositoryStore {
	return &RepositoryStore{
		log: printer.New(),
		k:   koanf.New("."),
	}
}

func (s *RepositoryStore) logArgs(args []string) []string {
	s.log.InfoF("%s %s", GitBinary, strings.Join(args, " "))
	return args
}

func (s *RepositoryStore) deleteFileHandler(file *os.File) {
	_ = file.Close()
	_ = os.Remove(file.Name())
}

func (s *RepositoryStore) FetchGitRepositories() ([]*domain.GitRepository, error) {
	if err := s.k.Load(file.Provider(ManagedReposFileName), yaml.Parser()); err != nil {
		return nil, err
	}
	var list []*domain.GitRepository
	var m []ManagedGitRepo
	if err := s.k.Unmarshal("repositories", &m); err != nil {
		return nil, err
	}
	gitBase := "git@github.com:"

	for _, repo := range m {
		u := parseUrl(repo, gitBase, s.DefaultNamespace)

		// TODO: Reimplement filtering

		gitUrl := domain.FromURL(u)
		root := s.toLocalFilePath(gitUrl.AsURL())
		gitRepository := domain.NewGitRepository(gitUrl, domain.NewFilePath(root))
		list = append(list, gitRepository)
	}
	return list, nil
}

func (s *RepositoryStore) toLocalFilePath(u *url.URL) string {
	p := strings.ReplaceAll(u.Path, "/", string(filepath.Separator))
	return filepath.Clean(filepath.Join(s.ParentDir, strings.ReplaceAll(u.Hostname(), ":", "-"), p))
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
