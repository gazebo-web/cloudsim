package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

func ScheduleTasks(cloudsim *platform.Platform, apps map[string]application.IApplication) {
	for _, app := range apps {
		tasks := app.RegisterTasks()
		for _, task := range tasks {
			cloudsim.Scheduler.DoAt(task.Job, task.Date)
		}
	}
}
