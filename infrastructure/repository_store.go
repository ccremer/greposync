package infrastructure

import (
	"errors"
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
	log printer.Printer
	k   *koanf.Koanf
	StoreConfig
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

func (s *RepositoryStore) Clone(repository *domain.GitRepository) error {
	if repository.RootDir.DirExists() {
		return errors.New("clone exists already")
	}
	dir := repository.RootDir.String()
	url := repository.URL

	// Don't want to expose credentials in the log, so we're not calling logArgs().
	s.log.InfoF("%s %s", GitBinary, strings.Join([]string{"clone", url.Redacted(), dir}, " "))

	out, stderr, err := execGitCommand(repository.RootDir, []string{"clone", url.String(), dir})
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.PrintF(out)
	return nil
}

func (s *RepositoryStore) Checkout(repository *domain.GitRepository) error {
	args := []string{"checkout"}
	if localExists, err := hasLocalBranch(repository, repository.CommitBranch); err != nil {
		return err
	} else if !localExists {
		// Checkout to new branch
		args = append(args, "-t", "-b")
	}
	args = append(args, repository.CommitBranch)

	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs(args))
	if err != nil {
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) Fetch(repository *domain.GitRepository) error {
	panic("implement me")
}

func (s *RepositoryStore) Reset(repository *domain.GitRepository) error {
	panic("implement me")
}

func (s *RepositoryStore) Pull(repository *domain.GitRepository) error {
	panic("implement me")
}

func (s *RepositoryStore) Commit(repository *domain.GitRepository, options domain.CommitOptions) error {
	f, err := os.CreateTemp("", "COMMIT_MSG_")
	if err != nil {
		return fmt.Errorf("failed to create temporary commit message file: %w", err)
	}
	defer s.deleteFileHandler(f)

	// Write commit message
	_, err = fmt.Fprint(f, options.Message)
	if err != nil {
		return err
	}

	args := []string{"commit", "-a", "-F", f.Name()}

	// Try to figure out if amend makes sense
	if options.Amend {
		if hasCommits, err := HasCommitsBetween(repository, repository.DefaultBranch, repository.CommitBranch); err != nil {
			return err
		} else if hasCommits {
			args = append(args, "--amend")
		}
	}

	// Commit
	out, stderr, err := execGitCommand(repository.RootDir, s.logArgs(args))
	if err != nil {
		s.log.InfoF(out)
		return mergeWithStdErr(err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *RepositoryStore) logArgs(args []string) []string {
	s.log.InfoF("%s %s", GitBinary, strings.Join(args, " "))
	return args
}

func (s *RepositoryStore) deleteFileHandler(file *os.File) {
	_ = file.Close()
	_ = os.Remove(file.Name())
}

func mergeWithStdErr(err error, stderr string) error {
	return fmt.Errorf("%w: %s", err, stderr)
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

		root := s.toLocalFilePath(u)
		domainRepo, err := domain.NewGitRepository(domain.FromURL(u), domain.NewFilePath(root), domain.LabelSet{})
		if err != nil {
			return list, err
		}
		list = append(list, domainRepo)
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
