package domain

import (
	"context"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

type CleanupService struct {
	instrumentation CleanupServiceInstrumentation
}

type CleanupPipeline struct {
	Repository *GitRepository
	ValueStore ValueStore

	files           []Path
	instrumentation CleanupServiceInstrumentation
}

func NewCleanupService(
	instrumentation CleanupServiceInstrumentation,
) *CleanupService {
	return &CleanupService{
		instrumentation: instrumentation,
	}
}

func (s *CleanupService) CleanupUnwantedFiles(pipe CleanupPipeline) error {
	pipe.instrumentation = s.instrumentation.WithRepository(pipe.Repository)
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("preflight check", pipe.preFlightCheck),
		pipeline.NewStepFromFunc("load files", pipe.loadFiles),
		pipeline.NewStepFromFunc("delete files", pipe.deleteFiles),
	).Run()
	return result.Err()
}

func (ctx *CleanupPipeline) preFlightCheck(_ context.Context) error {
	err := firstOf(
		checkIfArgumentNil(ctx.Repository, "Repository"),
		checkIfArgumentNil(ctx.ValueStore, "ValueStore"),
	)
	return err
}

func (ctx *CleanupPipeline) loadFiles(_ context.Context) error {
	files, err := ctx.ValueStore.FetchFilesToDelete(ctx.Repository)
	ctx.files = files
	return ctx.instrumentation.FetchedFilesToDelete(err, files)
}

func (ctx *CleanupPipeline) deleteFiles(_ context.Context) error {
	for _, file := range ctx.files {
		absoluteFile := ctx.Repository.RootDir.Join(file)
		if absoluteFile.FileExists() {
			if err := os.Remove(absoluteFile.String()); hasFailed(err) {
				return err
			}
			ctx.instrumentation.DeletedFile(absoluteFile)
		}
	}
	return nil
}
