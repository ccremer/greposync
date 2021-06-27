package repository

import "github.com/go-git/go-git/v5"

func (s *Service) PushToRemote() {

	if s.Config.SkipPush || s.Config.SkipCommit {
		Info("Skipping push")
		return
	}
	Info("git push")
	err := s.r.Push(&git.PushOptions{
		Force: s.Config.ForcePush,
	})
	CheckIfError(err)
}
