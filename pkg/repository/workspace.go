package repository

import (
	"errors"
	"fmt"
	"strings"
)

func (s *Repository) Clone() error {
	out, stderr, err := s.execGitCommand(s.logArgs("clone", s.GitConfig.Url.Redacted(), s.GitConfig.Dir)...)
	if err != nil {
		return fmt.Errorf("%w: %s", err, stderr)
	}
	s.log.PrintF(out)
	return nil
}

func (s *Repository) Checkout() error {
	args := []string{"checkout"}
	if localExists, err := s.localBranchExists(s.GitConfig.CommitBranch); err != nil {
		return err
	} else if !localExists {
		args = append(args, "-t", "-b")
	}
	args = append(args, s.GitConfig.CommitBranch)

	out, stderr, err := s.execGitCommand(s.logArgs(args...)...)
	if err != nil {
		return fmt.Errorf("%w: %s", err, stderr)
	}
	s.log.DebugF(out)
	return nil
}

func (s *Repository) Fetch() error {
	panic("implement me")
}

func (s *Repository) Reset() error {
	panic("implement me")
}

func (s *Repository) Pull() error {
	exists, err := s.remoteBranchExists(s.GitConfig.CommitBranch)
	if err != nil {
		return err
	}
	if exists {
		out, stderr, err := s.execGitCommand(s.logArgs("pull")...)
		if err != nil {
			return fmt.Errorf("%w: %s", err, stderr)
		}
		s.log.DebugF(out)
	}
	return nil
}

func (s *Repository) localBranchExists(branch string) (bool, error) {
	out, stderr, err := s.execGitCommand("branch", "--list")
	if err != nil {
		return false, errors.New(stderr)
	}
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Repository) remoteBranchExists(branch string) (bool, error) {
	out, stderr, err := s.execGitCommand("branch", "-r", "--list")
	if err != nil {
		return false, errors.New(stderr)
	}
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true, nil
		}
	}
	return false, nil
}
