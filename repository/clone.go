package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) CloneGitRepositoryAction() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("clone", s.Config.Url.Redacted(), s.Config.Dir)...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.PrintF(out)
		s.Config.DefaultBranch = s.GetDefaultBranch()
		return pipeline.Result{}
	}
}
