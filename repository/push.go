package repository

import "github.com/go-git/go-git/v5"

func PushToRemote(r *git.Repository) {

	Info("git push")
	// push using default options
	err := r.Push(&git.PushOptions{
		Force: true,
	})
	CheckIfError(err)
}
