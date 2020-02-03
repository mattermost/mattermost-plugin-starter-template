package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan"
)

func main() {
	// TODO: implement proper command line parameter parsing.
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "running: \n $ sync [plan.yaml] [plugin path]\n")
		os.Exit(1)
	}

	syncPlan, err := readPlan(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "coud not load plan: %s\n", err)
		os.Exit(1)
	}

	tplDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get current directory: %s\n", err)
		os.Exit(1)
	}

	pluginDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not determine plugin directory: %s\n", err)
		os.Exit(1)
	}

	tplRepo, err := plan.GetRepoContext(tplDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	pluginRepo, err := plan.GetRepoContext(pluginDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	planSetup := plan.Context{
		Template: tplRepo,
		Plugin:   pluginRepo,
	}
	err = syncPlan.Execute(planSetup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func readPlan(path string) (*plan.Plan, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file %q: %w", path, err)
	}

	var p plan.Plan
	err = yaml.Unmarshal(raw, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan yaml: %w", err)
	}

	return &p, err
}
