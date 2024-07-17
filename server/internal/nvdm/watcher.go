package nvdm

import (
	"context"
	"log"
	"time"

	"github.com/illbjorn/mm-nvd/nvd"
)

// CVEWatcher describes a daemon which will, on duration pollInterval, execute
// the nvd.CVEQuery. If any new results are identified (CVE.PublishDate >
// lastCheck) - NewCVEHandler handler will be called with the new CVE.
type CVEWatcher struct {
	// a cancellable context - the CVEWatcherGroup will call cancel() if an
	// Unregister() call occurs.
	ctx    context.Context
	cancel context.CancelFunc
	// The interval on which this worker will poll the NVD.
	pollInterval *time.Duration
	// The handler to be invoked when a new CVE is identified.
	handler NewCVEHandler
	// The assembled query to invoke and evaluate for new CVEs.
	cveQuery *nvd.CVEQuery
	// The last time this worker polled the NVD.
	lastCheck time.Time
}

func (w *CVEWatcher) CVEQuery() *nvd.CVEQuery {
	return w.cveQuery
}

func (w *CVEWatcher) Run() {
	// enter the worker loop
	for {
		select {
		case <-w.ctx.Done():
			return
		//case <-time.After(5 * time.Second):
		case <-time.After(*w.pollInterval):
			// perform the NVD lookup
			cves, err := w.cveQuery.Fetch()
			if err != nil {
				// TODO: we should log this somewhere useful.
				log.Println(err.Error())
			}

			// iterate CVEs
			for _, cve := range cves {
				// avoid an avalanche of messages on first-run
				//if w.lastCheck.Equal(time.Time{}) {
				//	break
				//}
				// only produce notifications for CVEs published after our last run
				if cve.Key.Published.After(w.lastCheck) && w.handler != nil {
					// invoke the handler
					w.handler(cve)
					break
				}
			}

			// update our last check time
			w.lastCheck = time.Now()
		}
	}
}
