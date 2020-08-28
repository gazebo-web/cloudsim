package actions

import (
	"errors"
	"strings"
)

var (
	// errNilJob is returned when trying to get the job index in an action for a nil job
	errNilJob = errors.New("cannot get job index for nil job")
	// errJobNotFound is returned when a job is not found in an action
	errJobNotFound = errors.New("job not found")
)

// Action contains a sequence of jobs, and performs a specific function.
// Actions are registered, launched and managed by action services.
// Action instances should only be used to define a sequence of states for registration in a service.
type Action struct {
	// Name contains the action name.
	// This field can be left empty as it will be filled in by a service when registering this action.
	Name string
	// Jobs contains the sequence of jobs processed to perform this action.
	Jobs Jobs
}

// NewAction creates a new Action containing a sequence of jobs.
func NewAction(jobs Jobs) (*Action, error) {
	// Create the action
	action := &Action{
		Jobs: jobs,
	}

	// Validate the sequence of jobs
	if err := action.Jobs.validate(); err != nil {
		return nil, err
	}

	// Register job data types
	action.Jobs.registerTypes(jobDataTypeRegistry)

	return action, nil
}

// getJobNames returns the sequence of job names for this action as a single comma-separated string.
func (a *Action) getJobNames() *string {
	jobNames := make([]string, len(a.Jobs))

	for i, job := range a.Jobs {
		jobNames[i] = job.Name
	}
	names := strings.Join(jobNames, ",")

	return &names
}

// getJobIndex gets the job index for a job.
func (a *Action) getJobIndex(jobName *string) (int, error) {
	// Sanity check
	if jobName == nil {
		return 0, errNilJob
	}

	for index, job := range a.Jobs {
		if job.Name == *jobName {
			return index, nil
		}
	}

	return 0, errJobNotFound
}
