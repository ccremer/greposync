package main

import (
	"log"

	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
)

func main() {
	dir := "git-repo-sync"
	repo, err := repository.CloneGitRepository("git@github.com:ccremer/git-repo-sync.git", dir)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{} {
		"Values": map[string]string {
			"name": dir,
		},
	}

	err = rendering.RenderTemplate(dir, data)
	if err != nil {
		log.Fatal(err)
	}

	err = repository.MakeCommit(repo)
	if err != nil {
		log.Fatal(err)
	}
}
