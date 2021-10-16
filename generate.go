//go:build generate
// +build generate

package main

// Run this file itself
//go:generate go run generate.go

// Generate dependency tree
//go:generate go run github.com/google/wire/cmd/wire

// Fix wire import
//go:generate sed -i -e "s|//+build|// +build|" wire_gen.go

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ccremer/greposync/cfg"
	"sigs.k8s.io/yaml"
)

func main() {
	createExampleConfig()
	createExampleLabelConfig()
}

func createExampleConfig() {
	exampleConfig := cfg.NewDefaultConfig()

	bytes, err := yaml.Marshal(exampleConfig)
	exit(err)
	writeFile(os.Getenv("REFERENCE_CONFIG_PATH"), bytes)

}

func createExampleLabelConfig() {
	config := map[string]interface{}{
		"repositoryLabels": cfg.RepositoryLabelMap{
			"greposync": {
				Name:        "greposync",
				Description: "updates from template repository",
				Color:       "#ededed",
				Delete:      false,
			},
		},
	}
	bytes, err := yaml.Marshal(config)
	exit(err)
	writeFile(os.Getenv("REFERENCE_LABELS_PATH"), bytes)
}

func writeFile(path string, bytes []byte) {
	exit(ioutil.WriteFile(path, bytes, 0775))
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
