package main

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

func main() {
	config := platform.NewConfig()
	cloudsim := platform.New(config)

	RegisterApplications(&cloudsim)
	RegisterRoutes(&cloudsim)

	if err := cloudsim.Start(context.Background()); err != nil {
		cloudsim.Logger.Critical("[CLOUDSIM] Error when starting platform up")
		cloudsim.Logger.Error(fmt.Sprintf("[ERROR] %v", err))
		for name, _ := range cloudsim.Applications {
			cloudsim.Logger.Info(fmt.Sprintf("\tRunning with application [%s]", name))
		}
		panic(err)
	}

	cloudsim.Server.Run()

	err := cloudsim.Stop(context.Background())
	if err != nil {
		cloudsim.Logger.Critical("[CLOUDSIM] Error on shutdown")
		cloudsim.Logger.Error(fmt.Sprintf("[ERROR] %v", err))
	}
	cloudsim.Transporter.Transport.Free()
}