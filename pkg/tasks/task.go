package tasks

import "time"

// Task represents a Job that needs to be run on a given Date.
// The scheduler uses a Task to schedule jobs for the applications.
type Task struct {
	Job func()
	Date time.Time
}
