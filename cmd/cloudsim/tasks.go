package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// ScheduleTasks gets all the tasks from each application and add them to the platform's scheduler.
func ScheduleTasks(cloudsim *platform.Platform, apps map[string]application.IApplication) {
	for _, app := range apps {
		tasks := app.RegisterTasks()
		for _, task := range tasks {
			cloudsim.Scheduler.DoAt(task.Job, task.Date)
		}
	}
}
