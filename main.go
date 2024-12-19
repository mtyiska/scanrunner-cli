package main

import (
	"fmt"
	"log"

	"github.com/mtyiska/scanrunner/cmd"
)

func main() {
	// Execute the CLI and handle any unexpected errors gracefully
	if err := recoverOnPanic(); err != nil {
		log.Fatalf("ScanRunner CLI encountered a fatal error: %v", err)
	}

	cmd.Execute()
}

// recoverOnPanic captures any unexpected panics to ensure graceful termination
func recoverOnPanic() (err error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			return e
		}
		return fmt.Errorf("%v", r)
	}
	return nil
}
