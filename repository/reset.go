package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) ResetRepository() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("fetch")...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		if out != "" {
			s.p.InfoF(out)
		}

		out, stderr, err = s.execGitCommand(s.logArgs("reset", "--hard")...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}
