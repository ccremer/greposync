package main

import (
	"log"

	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
)

func main() {
	dir := "git-repo-sync"
	repo := repository.PrepareWorkspace("git@github.com:ccremer/git-repo-sync.git", dir)

	data := map[string]interface{} {
		"Values": map[string]string {
			"name": dir,
		},
	}

	err := rendering.RenderTemplate(dir, data)
	if err != nil {
		log.Fatal(err)
	}

	repository.MakeCommit(repo)
	repository.ShowDiff(repo)
}
