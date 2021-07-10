package repository

import (
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

// Commit invokes git to stage all files and commit to the current branch.
func (s *Service) Commit() pipeline.ActionFunc {
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
			s.p.InfoF(out)
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}

// EnabledCommit returns true if commits are enabled.
func (s *Service) EnabledCommit() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return !s.Config.SkipCommit
	}
}

func (s *Service) Dirty() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		out, stderr, err := s.execGitCommand("status", "--short")
		if err != nil {
			s.p.WarnF("Could not determine working tree status: %s: %w", stderr, err)
			return true
		}
		if out == "" {
			s.p.InfoF("Nothing to commit, working tree clean")
			return false
		}
		return true
	}
}

func (s *Service) deleteFile(file *os.File) {
	_ = file.Close()
	err := os.Remove(file.Name())
	if err != nil {
		s.p.WarnF(err.Error())
	}
}
