package run_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/run"
)

// Test how the synchronization runner handles prechecks.
func TestSynchronizeRunsChecks(t *testing.T) {
	var calls []string

	// Function creating mock checkers that (optionally) fail.
	mockCheck := func(id string, fail bool) func() error {
		return func() error {
			calls = append(calls, id)
			if fail {
				return fmt.Errorf("%q failed", id)
			}
			return nil
		}

	}
	assert := assert.New(t)

	// Test with no failing runs.
	sync := run.Synchronize{
		Checks: []run.Check{
			mockCheck("1", false),
			mockCheck("2", false),
		},
	}

	err := sync.Run()
	assert.Nil(err)
	assert.Equal(calls, []string{"1", "2"})

	calls = []string{}
	// Test with a failing run.
	sync = run.Synchronize{
		Checks: []run.Check{
			mockCheck("1", false),
			mockCheck("2", true),
			mockCheck("3", false),
		},
	}
	err = sync.Run()
	assert.EqualError(err, "check failed: \"2\" failed")
	// Check after failure ("3") is not run.
	assert.Equal(calls, []string{"1", "2"})
}
