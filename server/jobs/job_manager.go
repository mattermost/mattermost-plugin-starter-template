package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/wiggin77/merror"

	"github.com/mattermost/mattermost-plugin-starter-template/server/config"
)

type Job interface {
	GetID() string
	OnConfigurationChange(cfg *config.Configuration) error
	Stop(timeout time.Duration) error
}

type JobManager struct {
	jobs   sync.Map
	logger loggerIface
}

func NewJobManager(logger loggerIface) *JobManager {
	return &JobManager{
		logger: logger,
	}
}

func (jm *JobManager) AddJob(job Job) error {
	_, loaded := jm.jobs.LoadOrStore(job.GetID(), job)

	if loaded {
		return fmt.Errorf("cannot add job: ID '%s' already exists", job.GetID())
	}
	return nil
}

func (jm *JobManager) RemoveJob(jobID string, timeout time.Duration) error {
	jobAny, loaded := jm.jobs.LoadAndDelete(jobID)
	if !loaded {
		return fmt.Errorf("cannot remove job: ID '%s' does not exist", jobID)
	}

	job := jobAny.(Job)
	if err := job.Stop(timeout); err != nil {
		jm.logger.Error("Error stopping job while removing from job manager", "err", err)
	}

	return nil
}

func (jm *JobManager) OnConfigurationChange(cfg *config.Configuration) error {
	merr := merror.New()

	jm.jobs.Range(func(_, v any) bool {
		job := v.(Job)
		if err := job.OnConfigurationChange(cfg); err != nil {
			merr.Append(err)
		}
		return true
	})
	return merr.ErrorOrNil()
}

func (jm *JobManager) Close(timeout time.Duration) error {
	merr := merror.New()

	jm.jobs.Range(func(_, v any) bool {
		job := v.(Job)
		if err := job.Stop(timeout); err != nil {
			merr.Append(err)
		}
		return true
	})
	return merr.ErrorOrNil()
}
