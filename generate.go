// +build generate

package main

// Run this file itself
//go:generate go run generate.go

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ccremer/greposync/cfg"
	"sigs.k8s.io/yaml"
)

func main() {
	config := cfg.NewDefaultConfig()

	bytes, err := yaml.Marshal(config)
	exit(err)
	exit(ioutil.WriteFile(os.Getenv("GODOC_YAML_DEFAULTS_PATH"), bytes, 0775))
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
