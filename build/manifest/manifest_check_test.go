package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManifestIsValid(t *testing.T) {
	// Save current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer os.Chdir(cwd) // restore after test

	// Change to repo root (where plugin.json exists)
	repoRoot := filepath.Join(cwd, "..", "..")
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("failed to change to repo root: %v", err)
	}

	// Set dummy build variables so Version can be populated
	BuildHashShort = "abc123"
	BuildTagCurrent = "v1.0.0"
	BuildTagLatest = "v1.0.0"

	manifest, err := findManifest()
	if err != nil {
		t.Fatalf("failed to find or parse manifest: %v", err)
	}

	if err := manifest.IsValid(); err != nil {
		t.Fatalf("manifest is invalid: %v", err)
	}
}
