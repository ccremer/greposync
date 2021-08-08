package gitrepo

import (
	"fmt"
	"os"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

type PrepareWorkspaceHandler struct {
	log       printer.Printer
	repoStore core.GitRepositoryStore
}

func NewPrepareWorkspaceHandler(repoStore core.GitRepositoryStore) *PrepareWorkspaceHandler {
	return &PrepareWorkspaceHandler{
		log:       printer.New(),
		repoStore: repoStore,
	}
}

func (s *PrepareWorkspaceHandler) Handle(source core.EventSource) core.EventResult {
	if source.Url == nil {
		return core.EventResult{Error: fmt.Errorf("no URL defined")}
	}
	repo, err := s.repoStore.FetchGitRepository(source.Url)
	if err != nil {
		return core.ToResult(source, err)
	}
	return core.ToResult(source, s.runPipeline(repo))
}

func (s *PrepareWorkspaceHandler) clone(ctx *pipelineContext) error {
	return ctx.repo.Clone()
}

func (s *PrepareWorkspaceHandler) fetch(ctx *pipelineContext) error {
	return ctx.repo.Fetch()
}

func (s *PrepareWorkspaceHandler) checkout(ctx *pipelineContext) error {
	return ctx.repo.Checkout()
}

func (s *PrepareWorkspaceHandler) reset(ctx *pipelineContext) error {
	return ctx.repo.Reset()
}

func (s *PrepareWorkspaceHandler) pull(ctx *pipelineContext) error {
	return ctx.repo.Pull()
}

func (s *PrepareWorkspaceHandler) dirExists(dir string) bool {
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return true
	}
	return false
}
