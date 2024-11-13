package jobs

import (
	"context"
	"fmt"

	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/mattermost/mattermost/server/public/pluginapi/cluster"
	"github.com/wiggin77/merror"

	"github.com/mattermost/mattermost-plugin-starter-template/server/config"
	"github.com/mattermost/mattermost-plugin-starter-template/server/store/kvstore"
)

type StarterTemplateJobSettings struct {
}

func (s *StarterTemplateJobSettings) Clone() *StarterTemplateJobSettings {
	return &StarterTemplateJobSettings{}
}

type StarterTemplateJob struct {
	mux      sync.Mutex
	settings *StarterTemplateJobSettings
	job      *cluster.Job
	runner   *runInstance

	id      string
	papi    plugin.API
	client  *pluginapi.Client
	kvstore kvstore.KVStore
}

func NewStarterTemplateJob(id string, api plugin.API, client *pluginapi.Client, kvstore kvstore.KVStore) (*StarterTemplateJob, error) {
	return &StarterTemplateJob{
		id:      id,
		papi:    api,
		client:  client,
		kvstore: kvstore,
	}, nil
}

func (j *StarterTemplateJob) GetID() string {
	return j.id
}

// OnConfigurationChange is called by the job manager whenenver the plugin settings have changed.
// It is suggested that you stop the current job (if any) and start a new job (if enabled) with new settings.
func (j *StarterTemplateJob) OnConfigurationChange(cfg *config.Configuration) error {
	j.client.Log.Debug("StarterTemplateJob: Configuration Changed")
	settings := &StarterTemplateJobSettings{}

	if err := j.Stop(time.Second * 10); err != nil {
		j.client.Log.Error("Error stopping Starter Template job for config change", "err", err)
	}

	j.client.Log.Debug("Preparing to start Starter Template job.")
	return j.start(settings)
}

// Start schedules a new job with specified settings.
func (j *StarterTemplateJob) start(settings *StarterTemplateJobSettings) error {
	j.mux.Lock()
	defer j.mux.Unlock()

	j.settings = settings

	job, err := cluster.Schedule(j.papi, j.id, j.nextWaitInterval, j.run)
	if err != nil {
		return fmt.Errorf("cannot start Starter Template job: %w", err)
	}
	j.job = job

	j.client.Log.Debug("Starter Template job started")

	return nil
}

// Stop stops the current job (if any). If the timeout is exceeded an error
// is returned.
func (j *StarterTemplateJob) Stop(timeout time.Duration) error {
	var job *cluster.Job
	var runner *runInstance

	j.mux.Lock()
	job = j.job
	runner = j.runner
	j.job = nil
	j.runner = nil
	j.mux.Unlock()

	merr := merror.New()

	if job != nil {
		if err := job.Close(); err != nil {
			merr.Append(fmt.Errorf("error closing job: %w", err))
		}
	}

	if runner != nil {
		if err := runner.stop(timeout); err != nil {
			merr.Append(fmt.Errorf("error stopping job runner: %w", err))
		}
	}

	j.client.Log.Debug("Starter Template Job stopped", "err", merr.ErrorOrNil())

	return merr.ErrorOrNil()
}

// func (j *StarterTemplateJob) getSettings() *StarterTemplateJobSettings {
// 	j.mux.Lock()
// 	defer j.mux.Unlock()
// 	return j.settings.Clone()
// }

// nextWaitInterval is called by the cluster job scheduler to determine how long to wait until the
// next job run.
func (j *StarterTemplateJob) nextWaitInterval(now time.Time, metaData cluster.JobMetadata) time.Duration {
	return time.Hour
}

func (j *StarterTemplateJob) RunFromAPI() {
	j.run()
}

func (j *StarterTemplateJob) run() {
	j.mux.Lock()
	oldRunner := j.runner
	j.mux.Unlock()

	if oldRunner != nil {
		j.client.Log.Error("Multiple Starter Template jobs scheduled concurrently; there can be only one")
		return
	}

	j.client.Log.Info("Running Starter Template Job")
	exitSignal := make(chan struct{})
	ctx, canceller := context.WithCancel(context.Background())

	runner := &runInstance{
		canceller:  canceller,
		exitSignal: exitSignal,
	}

	defer func() {
		canceller()
		close(exitSignal)

		j.mux.Lock()
		j.runner = nil
		j.mux.Unlock()
	}()

	var settings *StarterTemplateJobSettings
	j.mux.Lock()
	j.runner = runner
	settings = j.settings.Clone()
	j.mux.Unlock()

	// Job logic goes here.

	_ = ctx
	_ = settings
}
