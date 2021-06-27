package main

import (
	"log"
	"path"

	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
)

func main() {
	dir := "git-repo-sync"
	repo := repository.PrepareWorkspace("git@github.com:ccremer/git-repo-sync.git", dir)

	rendering.LoadConfigFile("config_defaults.yml")
	syncFile := path.Join("repos",dir, ".sync.yml")
	rendering.LoadConfigFile(syncFile)

	data := map[string]interface{}{
		"Values": rendering.Unmarshal("README.md/test"),
	}

	err := rendering.RenderTemplate(dir, data)
	if err != nil {
		log.Fatal(err)
	}

	//repository.MakeCommit(repo)
	repository.ShowDiff(repo)
	//repository.PushToRemote(repo)
	//repository.CreatePR(repo)
}
