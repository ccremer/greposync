package rendering

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
)

// DeleteUnwantedFiles goes through the sync config and deletes files from repositories that aren't targeted by the templates.
// The config need to have the `delete` flag set to true.
// Only files are deleted, not directories.
func (r *Renderer) DeleteUnwantedFiles() pipeline.ActionFunc {
	return func() pipeline.Result {
		files := r.searchOrphanedFiles()
		for _, relativePath := range files {
			targetPath := path.Clean(path.Join(r.cfg.Git.Dir, relativePath))
			err := r.deleteFileIfExists(targetPath)
			if err != nil {
				return pipeline.Result{Err: err}
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) searchOrphanedFiles() []string {
	filePaths := make([]string, 0)
	allKeys := r.k.Raw()
	// Go through all top-level keys, which are the file names
	for filePath, values := range allKeys {
		// If the filename is already handled by the template renderer, ignore it.
		// Otherwise, add files that have deletion flag, but ignore directories
		if _, found := r.parser.templates[filePath]; !found && pathIsFile(filePath) {
			if val, ok := values.(map[string]interface{}); ok {
				if filePath == ":globals" {
					// can't delete file named ':globals' anyway
					continue
				}
				if val["delete"] == true {
					filePaths = append(filePaths, filePath)
				}
			}
		}
	}
	return filePaths
}

func pathIsFile(filePath string) bool {
	return !strings.HasSuffix(filePath, string(filepath.Separator))
}

func (r *Renderer) deleteFileIfExists(targetPath string) error {
	if fileExists(targetPath) {
		r.p.InfoF("Deleting file due to 'delete' flag being set: %s", targetPath)
		return os.Remove(targetPath)
	}
	return nil
}
