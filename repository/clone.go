package repository

import (
	"os"
)

func (s *Service) PrepareWorkspace() {
	if _, err := os.Stat(s.Config.Dir); os.IsNotExist(err) {
		s.CloneGitRepository()
		s.SwitchBranch()
		return
	}
	s.ResetRepository()
	s.SwitchBranch()
	s.Pull()
}

func (s *Service) CloneGitRepository() {
	out, _, err := s.execGitCommand(s.logArgs("clone", s.Config.Url.Redacted(), s.Config.Dir)...)
	s.p.CheckIfError(err)
	s.p.PrintF(out)
	s.Config.DefaultBranch = s.GetDefaultBranch()
}
