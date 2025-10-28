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

	// Restore working directory after test
	defer func() {
		if chdirErr := os.Chdir(cwd); chdirErr != nil {
			t.Logf("warning: failed to restore working directory: %v", chdirErr)
		}
	}()

	// Change to repo root (where plugin.json exists)
	repoRoot := filepath.Join(cwd, "..", "..")
	err = os.Chdir(repoRoot)
	if err != nil {
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

	if valErr := manifest.IsValid(); valErr != nil {
		t.Fatalf("manifest is invalid: %v", valErr)
	}
}

func TestManifestIsInvalid(t *testing.T) {
	// Save current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Restore working directory after test
	defer func() {
		if chdirErr := os.Chdir(cwd); chdirErr != nil {
			t.Logf("warning: failed to restore working directory: %v", chdirErr)
		}
	}()

	// Change to repo root (where plugin.json exists)
	repoRoot := filepath.Join(cwd, "..", "..")
	err = os.Chdir(repoRoot)
	if err != nil {
		t.Fatalf("failed to change to repo root: %v", err)
	}

	// Set dummy build variables
	BuildHashShort = "abc123"
	BuildTagCurrent = "v1.0.0"
	BuildTagLatest = "v1.0.0"

	// Load the manifest
	manifest, err := findManifest()
	if err != nil {
		t.Fatalf("failed to find or parse manifest: %v", err)
	}

	// Invalidate the manifest
	manifest.Id = ""

	if valErr := manifest.IsValid(); valErr == nil {
		t.Fatal("expected manifest.IsValid() to fail when Id is empty, but it passed")
	}
}
