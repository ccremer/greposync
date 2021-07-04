package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

// ResetRepository invokes git to reset the git repository to discard local changes.
func (s *Service) ResetRepository() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("reset", "--hard")...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}

// Fetch invokes git to fetch remote references.
func (s *Service) Fetch() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("fetch")...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		if out != "" {
			s.p.InfoF(out)
		}
		return pipeline.Result{}
	}
}
