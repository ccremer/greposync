package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

// CloneGitRepository invokes git to clone the repository.
func (s *Service) CloneGitRepository() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("clone", s.Config.Url.Redacted(), s.Config.Dir)...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.PrintF(out)
		return pipeline.Result{}
	}
}
