package main

import (
	"time"

	"github.com/ccremer/greposync/cli"
)

var (
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")
)

func main() {
	cli.CreateCLI(version, commit, date)
	cli.Run()
}
