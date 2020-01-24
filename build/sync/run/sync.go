package run

import (
	"fmt"
)

// Synchronize defines the set and order of operations to be executed
// to synchronize a plugin repository with the template.
type Synchronize struct {
	// Checks are run before performing the sync.
	Checks []Check
}

// Check returns a non-nil error if the check fails.
type Check func() error

func (s Synchronize) Run() error {
	for _, check := range s.Checks {
		if err := check(); err != nil {
			return fmt.Errorf("check failed: %w", err)
		}
	}
	return nil
}
