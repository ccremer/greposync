package repository

import (
	"fmt"
	"log"
)

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
