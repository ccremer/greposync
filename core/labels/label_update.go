package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

// LabelUpdateEvent identifies an event that will update labels in a remote Git repository.
const LabelUpdateEvent core.EventName = "core:label-update"

// LabelUpdateHandler contains the business logic to interact with labels on supported core.GitHostingProvider.
type LabelUpdateHandler struct {
	repoStore core.GitRepositoryStore
}

type pipelineContext struct {
	repo core.GitRepository
	log  printer.Printer
}

// NewLabelUpdateHandler returns a new core LabelUpdateHandler instance.
func NewLabelUpdateHandler(repoProvider core.GitRepositoryStore) *LabelUpdateHandler {
	return &LabelUpdateHandler{
		repoStore: repoProvider,
	}
}

// Handle implements core.EventHandler.
func (s *LabelUpdateHandler) Handle(source core.EventSource) core.EventResult {
	if source.Url == nil {
		return core.EventResult{Error: fmt.Errorf("no URL defined")}
	}
	repo, err := s.repoStore.FetchGitRepository(source.Url)
	if err != nil {
		return core.ToResult(source, err)
	}
	return core.ToResult(source, s.runPipeline(repo))
}

func (s *LabelUpdateHandler) runPipeline(repo core.GitRepository) error {
	ctx := &pipelineContext{
		repo: repo,
		log:  printer.New().SetName(repo.GetConfig().URL.GetRepositoryName()),
	}

	logger := printer.PipelineLogger{Logger: ctx.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("edit labels", toAction(ctx, s.createOrUpdateLabels)),
		pipeline.NewStep("delete labels", toAction(ctx, s.deleteLabels)),
	).Run()
	return result.Err
}

func toAction(ctx *pipelineContext, action func(ctx *pipelineContext) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action(ctx)}
	}
}

func (s *LabelUpdateHandler) createOrUpdateLabels(ctx *pipelineContext) error {
	labels := ctx.repo.GetLabels()
	labels = filterActiveLabels(labels)
	if len(labels) <= 0 {
		return nil
	}
	for _, label := range labels {
		changed, err := label.Ensure()
		if err != nil {
			return err
		}
		if changed {
			ctx.log.InfoF("Label '%s' changed", label.GetName())
		} else {
			ctx.log.InfoF("Label '%s' unchanged", label.GetName())
		}
	}
	return nil
}

func filterActiveLabels(labels []core.Label) []core.Label {
	filtered := make([]core.Label, 0)
	for _, label := range labels {
		if !label.IsInactive() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func (s *LabelUpdateHandler) deleteLabels(ctx *pipelineContext) error {
	labels := ctx.repo.GetLabels()
	labels = filterDeadLabels(labels)
	if len(labels) <= 0 {
		return nil
	}
	for _, label := range labels {
		deleted, err := label.Delete()
		if err != nil {
			return err
		}
		if deleted {
			ctx.log.InfoF("Label '%s' deleted", label.GetName())
		} else {
			ctx.log.InfoF("Label '%s' not deleted (not existing)", label.GetName())
		}
	}
	return nil
}

func filterDeadLabels(labels []core.Label) []core.Label {
	var filtered []core.Label
	for _, label := range labels {
		if label.IsInactive() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}
