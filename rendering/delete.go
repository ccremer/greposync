package rendering

import (
	"os"
	"path"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/core"
)

// DeleteUnwantedFiles goes through the sync config and deletes files from repositories that aren't targeted by the templates.
// The config need to have the `delete` flag set to true.
// Only files are deleted, not directories.
func (r *Renderer) DeleteUnwantedFiles() pipeline.ActionFunc {
	return func() pipeline.Result {
		files, err := r.valueStore.FetchFilesToDelete(&core.GitRepositoryProperties{
			URL:     core.FromURL(r.cfg.Git.Url),
			RootDir: r.cfg.Git.Dir,
		})
		if err != nil {
			return pipeline.Result{Err: err}
		}
		for _, file := range files {
			err := r.deleteFileIfExists(path.Join(r.cfg.Git.Dir, file))
			if err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) deleteFileIfExists(targetPath string) error {
	if fileExists(targetPath) {
		r.p.InfoF("Deleting file due to 'delete' flag being set: %s", targetPath)
		return os.Remove(targetPath)
	}
	return nil
}
