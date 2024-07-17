package nvdm

import (
	"context"
	"time"

	"github.com/illbjorn/mm-nvd/nvd"
)

// NewCVEHandler describes a callback function which can be passed to
// CVEWatcherGroup.Register() to be called when a new CVE matching that worker's
// nvd.CVEQuery has been identified.
type NewCVEHandler func(cve nvd.CVE)

// NewCVEWatcherGroup initializes a CVEWatcherGroup with some default values.
func NewCVEWatcherGroup(ctx context.Context) *CVEWatcherGroup {
	return &CVEWatcherGroup{
		ctx:         ctx,
		pollingRate: 10 * time.Second,
		workers:     make(map[string]*CVEWatcher, 64),
	}
}

// CVEWatcherGroup describes a CVEWatcher factory.
type CVEWatcherGroup struct {
	ctx         context.Context
	pollingRate time.Duration
	workers     map[string]*CVEWatcher
}

// SetPollRate adjusts the rate at which all workers will query the NVD.
func (wg *CVEWatcherGroup) SetPollRate(rate time.Duration) {
	wg.pollingRate = rate
}

// Register registers a new CVEQuery worker with NewCVEHandler cb.
//
// The key parameter is used to associate a unique ID to an individual worker
// and prevent duplicate worker creation.
func (wg *CVEWatcherGroup) Register(
	key string,
	query *nvd.CVEQuery,
	cb NewCVEHandler,
) {
	// apply some query defaults
	query.ResultsPerPage(100).
		PublishedWithin(3 * 24 * 60 * time.Minute)

	// create the child context
	cctx, cancel := context.WithCancel(wg.ctx)

	// initialize the worker
	w := &CVEWatcher{
		ctx:          cctx,
		cancel:       cancel,
		pollInterval: &wg.pollingRate,
		handler:      cb,
		cveQuery:     query,
		lastCheck:    time.Now(),
	}

	// begin the worker
	go w.Run()

	// hang onto the worker for cancellation if necessary
	wg.workers[key] = w
}

// Unregister calls the context.CancelFunc to terminate the worker associated
// with unique ID 'key', if it exists.
func (wg *CVEWatcherGroup) Unregister(key string) error {
	// locate the worker
	worker, ok := wg.workers[key]
	if !ok {
		return ErrWorkerIDNotFound
	}

	// cancel the worker
	worker.cancel()

	return nil
}

// Worker returns a pointer to a worker associated with unique ID 'key', if it
// exists.
func (wg *CVEWatcherGroup) Worker(key string) *CVEWatcher {
	if worker, ok := wg.workers[key]; ok {
		return worker
	}

	return nil
}
