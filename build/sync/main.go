package main

import (
	"fmt"
	"os"
	//	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		reportError(fmt.Errorf("failed to get current directory: %w", err))
	}
	println(wd)
}

func reportError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
