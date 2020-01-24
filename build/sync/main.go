package main

import (
	"fmt"
	"os"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/git"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		reportError(fmt.Errorf("failed to get current directory: %w", err))
	}

	repo, err := git.Open(wd)
	if err != nil {
		reportError(err)
	}

	clean, err := repo.IsClean()
	if err != nil {
		reportError(err)
	}
	if !clean {
		reportError(fmt.Errorf("template repository is not clean"))
	}
}

func reportError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}
