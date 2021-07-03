package repository

import (
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) MakeCommit() pipeline.ActionFunc {
	return func() pipeline.Result {
		f, err := os.CreateTemp("", "COMMIT_MSG_")
		if err != nil {
			return pipeline.Result{Err: err}
		}
		defer s.deleteFile(f)

		// Write commit message
		_, err = fmt.Fprint(f, s.Config.CommitMessage)
		if err != nil {
			return pipeline.Result{Err: err}
		}

		// Commit
		args := []string{"commit", "-a", "-F", f.Name()}
		if s.Config.Amend {
			args = append(args, "--amend")
		}

		out, stderr, err := s.execGitCommand(s.logArgs(args...)...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}

func (s *Service) SkipCommit() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return s.Config.SkipCommit
	}
}

func (s *Service) deleteFile(file *os.File) {
	_ = file.Close()
	err := os.Remove(file.Name())
	if err != nil {
		s.p.WarnF(err.Error())
	}
}
