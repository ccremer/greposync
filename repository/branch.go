package repository

import (
	"strings"
)

func (s *Service) SwitchBranch() {
	if s.Config.SkipReset || s.Config.CommitBranch == "" {
		return
	}
	s.CheckoutBranch(s.Config.CommitBranch)
}

func (s *Service) CheckoutBranch(branch string) {
	out, _, err := s.execGitCommand(s.logArgs("checkout", branch)...)
	s.p.CheckIfError(err)
	s.p.DebugF(out)
}

func (s *Service) GetDefaultBranch() string {
	out, _, err := s.execGitCommand("remote", "show", "origin")
	s.p.CheckIfError(err)
	lines := strings.Split(out, "\n")
	head := "HEAD branch: "
	for _, line := range lines {
		str := strings.TrimSpace(line)
		if strings.Contains(str, head) {
			return strings.TrimPrefix(str, head)
		}
	}
	s.p.WarnF("No default branch detected. Fall back to master")
	return "master"
}

func (s *Service) localBranchExists(branch string) bool {
	out, _, err := s.execGitCommand("branch", "--list")
	s.p.CheckIfError(err)
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true
		}
	}
	return false
}

func (s *Service) remoteBranchExists(branch string) bool {
	out, _, err := s.execGitCommand("branch", "-r", "--list")
	s.p.CheckIfError(err)
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true
		}
	}
	return false
}
