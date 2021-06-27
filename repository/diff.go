package repository

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func ShowDiff(r *git.Repository) {

	// Getting the latest commit on the current branch
	Info("git log -1")

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)
	fmt.Println(commit)

	Info("retrieve parent commit")
	parent, err := commit.Parent(0)


	Info("getting previous commit %s", parent)
	patch, err := parent.Patch(commit)
	CheckIfError(err)

	prettyPrint(patch.String())
}

func prettyPrint(diff string) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			fmt.Fprintf(os.Stdout, "\x1b[31;1m%s\x1b[0m\n", line)
			continue
		}
		if strings.HasPrefix(line, "+") {
			fmt.Fprintf(os.Stdout, "\x1b[32;1m%s\x1b[0m\n", line)
			continue
		}
		fmt.Println(line)
	}
}
