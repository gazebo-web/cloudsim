package main

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

func main() {
	var cloudsim *platform.Platform
	var applications map[string]application.IApplication

	config := platform.NewConfig()
	cloudsim = platform.New(config)

	if err := cloudsim.Start(cloudsim.Context); err != nil {
		cloudsim.Logger.Critical(fmt.Sprintf("[CLOUDSIM|CRITICAL] Error when initializing cloudsim\n%v", err))
		for name, _ := range applications {
			cloudsim.Logger.Info(fmt.Sprintf("\tRunning with application [%s]", name))
		}
		panic(err)
	}

	RegisterApplications(cloudsim, &applications)
	RegisterMonitors(applications)
	RebuildState(cloudsim, applications)
	RegisterRoutes(cloudsim, applications)
	ScheduleTasks(cloudsim, applications)

	cloudsim.Server.Run()

	ShutdownApplications(applications)

	err := cloudsim.Stop(cloudsim.Context)
	if err != nil {
		cloudsim.Logger.Critical(fmt.Sprintf("[CLOUDSIM|CRITICAL] Error on shutdown\n%v", err))
	}
	cloudsim.Transport.Stop()
}
