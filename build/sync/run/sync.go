// Package run defines the logic for running updates
// with the template repository.
package run

import (
	"fmt"
)

// Synchronize defines the set and order of operations to be executed
// to synchronize a plugin repository with the template.
type Synchronize struct {
	// Checks are run before performing the sync.
	Checks []Check
	// Updates are paths that need to be updated.
	Updates map[string]Update
}

// Check returns a non-nil error if the check fails.
type Check func() error

// Update runs the update operation.
type Update func() error

/*
updates = Overwrite(Recreate(true))
updates = OneOf(
     Overwrite(IfUnaltered()),
     MergeGoMod())



*/

func (s Synchronize) Run() error {
	for _, check := range s.Checks {
		if err := check(); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}
	}
	return nil
}
