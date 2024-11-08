package jobs

import (
	"fmt"
	"time"
)

type loggerIface interface {
	// Error logs an error message, optionally structured with alternating key, value parameters.
	Error(message string, keyValuePairs ...interface{})

	// Warn logs an error message, optionally structured with alternating key, value parameters.
	Warn(message string, keyValuePairs ...interface{})

	// Info logs an error message, optionally structured with alternating key, value parameters.
	Info(message string, keyValuePairs ...interface{})

	// Debug logs an error message, optionally structured with alternating key, value parameters.
	Debug(message string, keyValuePairs ...interface{})
}

type runInstance struct {
	canceller  func()        // called to stop a currently executing run
	exitSignal chan struct{} // closed when the currently executing run has exited
}

func (r *runInstance) stop(timeout time.Duration) error {
	// cancel the run
	r.canceller()

	// wait for it to exit
	select {
	case <-r.exitSignal:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("waiting on job to stop timed out after %s", timeout.String())
	}
}
