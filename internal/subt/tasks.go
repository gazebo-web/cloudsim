package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"time"
)

// RegisterTasks returns an array of the tasks that need to be executed by the scheduler.
func (app *SubT) RegisterTasks() []tasks.Task {
	return []tasks.Task{
		{
			Job:  func() {

			},
			Date: time.Time{},
		},
	}
}