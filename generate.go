// +build generate

package main

// Run this file itself
//go:generate go run generate.go

// Generate fakes
//go:generate go generate ./core

// Generate dependency tree
//go:generate go run github.com/google/wire/cmd/wire

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ccremer/greposync/cfg"
	"sigs.k8s.io/yaml"
)

func main() {
	createExampleConfig()
}

func createExampleConfig() {
	exampleConfig := cfg.NewDefaultConfig()

	bytes, err := yaml.Marshal(exampleConfig)
	exit(err)
	writeFile(os.Getenv("GODOC_YAML_DEFAULTS_PATH"), bytes)
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
