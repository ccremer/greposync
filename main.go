package main

import (
	"log"
	"path"

	"github.com/ccremer/git-repo-sync/printer"
	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
)

func main() {
	printer.DefaultPrinter.SetLevel(printer.LevelDebug)

	services := repository.NewServicesFromFile("managed_repos.yml", "repos", "ccremer")

	for _, repoService := range services {
		repoService.PrepareWorkspace()

		rendering.LoadConfigFile("config_defaults.yml")
		syncFile := path.Join(repoService.Config.GitDir, ".sync.yml")
		rendering.LoadConfigFile(syncFile)

		data := map[string]interface{}{
			"Values": rendering.Unmarshal("README.md/test"),
		}

		err := rendering.RenderTemplate(repoService.Config.GitDir, data)
		if err != nil {
			log.Fatal(err)
		}

		repoService.MakeCommit()
		repoService.ShowDiff()
		repoService.PushToRemote()
		repoService.CreatePR()
	}
}
