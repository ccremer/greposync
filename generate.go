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
	"strings"

	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/cfg"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"sigs.k8s.io/yaml"
)

func main() {
	createExampleConfig()
	createExampleLabelConfig()
}

func createExampleConfig() {
	exampleConfig := generateStructure(
		flags.NewPRBodyFlag(nil),
		flags.NewPRCreateFlag(nil),
		flags.NewPRSubjectFlag(nil),
		flags.NewPRTargetBranchFlag(nil),
		flags.NewPRLabelsFlag(nil),

		flags.NewGitRootDirFlag(nil),
		flags.NewGitCommitMessageFlag(nil),
		flags.NewGitCommitBranchFlag(nil),
		flags.NewGitDefaultNamespaceFlag(nil),
		flags.NewGitForcePushFlag(nil),

		flags.NewShowDiffFlag(nil),
		flags.NewShowLogFlag(nil),

		flags.NewTemplateRootDirFlag(nil),
	)

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

func generateStructure(flags ...cli.Flag) map[string]interface{} {
	res := make(map[string]interface{}, 0)
	for _, flag := range flags {
		paths := strings.Split(flag.Names()[0], ".")
		category := ""
		name := flag.Names()[0]
		var value interface{}
		if len(paths) > 1 {
			category = paths[0]
			name = paths[1]
		}

		if strFlag, ok := flag.(*altsrc.StringFlag); ok {
			value = strFlag.Value
		}
		if intFlag, ok := flag.(*altsrc.IntFlag); ok {
			value = intFlag.Value
		}
		if pathFlag, ok := flag.(*altsrc.PathFlag); ok {
			value = pathFlag.Value
		}
		if boolFlag, ok := flag.(*altsrc.BoolFlag); ok {
			value = boolFlag.Value
		}
		if strSliceFlag, ok := flag.(*altsrc.StringSliceFlag); ok {
			slice := strSliceFlag.Value.Value()
			if slice == nil {
				slice = []string{}
			}
			value = slice
		}
		if value == nil {
			exit(fmt.Errorf("unrecognized flag type: %s", flag.Names()[0]))
		}

		if category == "" {
			res[name] = value
		} else {
			if res[category] == nil {
				res[category] = map[string]interface{}{}
			}
			res[category].(map[string]interface{})[name] = value
		}
	}
	return res
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
