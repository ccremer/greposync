package repository

import (
	"github.com/go-git/go-git/v5"
)

func (s *Service) PushToRemote() {

	if s.Config.SkipPush {
		s.p.WarnF("Skipped: push")
		return
	}
	s.p.InfoF("git push")
	err := s.r.Push(&git.PushOptions{
		Force: s.Config.ForcePush,
	})
	s.p.CheckIfError(err)
}
