package repository

import (
	"fmt"
	"os"
)

func (s *Service) MakeCommit() {
	if s.Config.SkipCommit {
		s.p.WarnF("Skipped: git commit")
		return
	}
	f, err := os.CreateTemp("", "COMMIT_MSG_")
	s.p.CheckIfError(err)
	defer s.deleteFile(f)
	args := []string{"commit", "-a", "-F", f.Name()}
	if s.Config.Amend {
		args = append(args, "--amend")
	}

	_, err = fmt.Fprint(f, s.Config.CommitMessage)
	s.p.CheckIfError(err)
	out, _, err := s.execGitCommand(s.logArgs(args...)...)
	s.p.DebugF(out)
	s.p.CheckIfError(err)
}

func (s *Service) deleteFile(file *os.File) {
	_ = file.Close()
	err := os.Remove(file.Name())
	if err != nil {
		s.p.WarnF(err.Error())
	}
}
