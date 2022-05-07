package main

var (
	version = "unknown"
	commit  = "-dirty-"
	date    = ""
)

func main() {
	injector := initInjector()
	injector.RunApp()
}
