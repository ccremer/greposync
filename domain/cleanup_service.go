package domain

import (
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

type CleanupService struct{}

type CleanupContext struct {
	Repository *GitRepository
	ValueStore ValueStore

	files []Path
}

func NewCleanupService() *CleanupService {
	return &CleanupService{}
}

func (s *CleanupService) CleanupUnwantedFiles(ctx CleanupContext) error {

	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStep("preflight check", ctx.preFlightCheck()),
		pipeline.NewStep("load files", ctx.toAction(ctx.loadFiles)),
		pipeline.NewStep("delete files", ctx.toAction(ctx.deleteFiles)),
	).Run()
	return result.Err
}

func (ctx *CleanupContext) preFlightCheck() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		err := firstOf(
			checkIfArgumentNil(ctx.Repository, "Repository"),
			checkIfArgumentNil(ctx.ValueStore, "ValueStore"),
		)
		return pipeline.Result{Err: err}
	}
}

func (ctx *CleanupContext) toAction(action func() error) pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		return pipeline.Result{Err: action()}
	}
}

func (ctx *CleanupContext) loadFiles() error {
	files, err := ctx.ValueStore.FetchFilesToDelete(ctx.Repository)
	ctx.files = files
	return err
}

func (ctx *CleanupContext) deleteFiles() error {
	for _, file := range ctx.files {
		absoluteFile := ctx.Repository.RootDir.Join(file)
		if absoluteFile.FileExists() {
			if err := os.Remove(absoluteFile.String()); hasFailed(err) {
				return err
			}
		}
	}
	return nil
}
