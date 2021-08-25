package repositorystore

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	giturls "github.com/whilp/git-urls"
)

type RepositoryStore struct {
	StoreConfig
	log        printer.Printer
	k          *koanf.Koanf
	prStore    domain.PullRequestStore
	labelStore domain.LabelStore
}

// ManagedGitRepo is the representation of the managed git repos in the config file.
type ManagedGitRepo struct {
	Name string
}

type StoreConfig struct {
	ParentDir        string
	DefaultNamespace string
}

type pipelineContext struct {
	url    *domain.GitURL
	labels domain.LabelSet
	pr     *domain.PullRequest
	repo   *domain.GitRepository
}

var (
	// ManagedReposFileName is the base file name where managed git repositories config is searched.
	ManagedReposFileName = "managed_repos.yml"
)

func NewRepositoryStore(prStore domain.PullRequestStore, labelStore domain.LabelStore) *RepositoryStore {
	return &RepositoryStore{
		log:        printer.New(),
		k:          koanf.New("."),
		labelStore: labelStore,
		prStore:    prStore,
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

		gitUrl := domain.FromURL(u)
		ctx := &pipelineContext{
			url: gitUrl,
		}
		p := pipeline.NewPipeline().WithSteps(
			pipeline.NewStep("fetch labels", ctx.fetchLabels(s.labelStore)),
			pipeline.NewStep("fetch PR", ctx.findPR(s.prStore)),
			pipeline.NewStep("create repo", ctx.createRepo(s)),
		)
		result := p.Run()
		if result.IsFailed() {
			return nil, result.Err
		}
		list = append(list, ctx.repo)
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

func (ctx *pipelineContext) fetchLabels(labelStore domain.LabelStore) pipeline.ActionFunc {
	return func() pipeline.Result {
		labels, err := labelStore.FetchLabelsForRepository(ctx.url)
		ctx.labels = labels
		return pipeline.Result{Err: err}
	}
}

func (ctx *pipelineContext) findPR(prStore domain.PullRequestStore) pipeline.ActionFunc {
	return func() pipeline.Result {
		pr, err := prStore.FindMatchingPullRequest(ctx.repo)
		ctx.pr = pr
		return pipeline.Result{Err: err}
	}
}

func (ctx *pipelineContext) createRepo(s *RepositoryStore) pipeline.ActionFunc {
	return func() pipeline.Result {
		root := s.toLocalFilePath(ctx.url.AsURL())
		domainRepo, err := domain.NewGitRepository(ctx.url, domain.NewFilePath(root), ctx.labels)
		ctx.repo = domainRepo
		return pipeline.Result{Err: err}
	}
}
