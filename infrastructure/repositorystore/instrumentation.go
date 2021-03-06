package repositorystore

import (
	"fmt"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

type RepositoryStoreInstrumentation struct {
	log logr.Logger
}

func NewRepositoryStoreInstrumentation(factory logging.LoggerFactory) *RepositoryStoreInstrumentation {
	return &RepositoryStoreInstrumentation{
		log: factory.NewGenericLogger(""),
	}
}

func (i *RepositoryStoreInstrumentation) attemptCloning(repository *domain.GitRepository) {
	// Don't want to expose credentials in the log, so we're not calling logArgs().
	i.log.WithName(repository.URL.GetFullName()).Info(fmt.Sprintf("%s %s", GitBinary, strings.Join([]string{"clone", repository.URL.Redacted(), repository.RootDir.String()}, " ")))
}

func (i *RepositoryStoreInstrumentation) logInfo(repository *domain.GitRepository, line string) {
	if line != "" {
		i.log.WithName(repository.URL.GetFullName()).Info(line)
	}
}

func (i *RepositoryStoreInstrumentation) logGitArguments(repository *domain.GitRepository, level int, args []string) []string {
	i.log.WithName(repository.URL.GetFullName()).V(level).Info(fmt.Sprintf("%s %s", GitBinary, strings.Join(args, " ")))
	return args
}

func (i *RepositoryStoreInstrumentation) logDebugInfo(repository *domain.GitRepository, line string) {
	if line != "" {
		i.log.WithName(repository.URL.GetFullName()).V(1).Info(line)
	}
}

func (i *RepositoryStoreInstrumentation) skipRepository(url *domain.GitURL) {
	i.log.Info("Skipping repository due to filters", "url", url.GetFullName())
}

func (i *RepositoryStoreInstrumentation) loadRepositoryConfigFile(name string) {
	i.log.V(1).Info("Loading config file", "name", name)
}
