package main

import (
	"fmt"
	"os"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/run"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		reportError(fmt.Errorf("failed to get current directory: %w", err))
	}

	sync := run.Synchronize{
		Checks: []run.Check{
			run.RepoIsClean(wd),
		},
	}
	err = sync.Run()
	if err != nil {
		reportError(err)
	}
}

func reportError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
