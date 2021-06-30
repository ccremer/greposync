package rendering

import (
	"path"
)

func (r *Renderer) ConstructMetadata() Values {
	d := Values{
		"Repository":  r.cfg.Git,
		"PullRequest": r.cfg.PullRequest,
	}
	return d
}

func (r *Renderer) ConstructTemplateMetadata(targetPath string) Values {
	return Values{
		"Metadata": Values{
			"Path":     targetPath,
			"FileName": path.Base(targetPath),
		},
	}
}

func (r *Renderer) GetMetadata(targetPath string) Values {
	return map[string]interface{}{
		"Path":     targetPath,
		"FileName": path.Base(targetPath),
		"RepoName": path.Base(r.cfg.Template.RootDir),
	}
}
