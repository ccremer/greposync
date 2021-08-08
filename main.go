package main

import (
	"time"
)

var (
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")
)

func main() {
	injector := initInjector()
	injector.RegisterHandlers()
	injector.RunApp()
}
