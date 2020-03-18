package main

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

func main() {
	config := platform.NewConfig()
	cloudsim := platform.New(config)
	RegisterRoutes(cloudsim)
	cloudsim.Server.Run()
	cloudsim.Stop(context.Background())
	cloudsim.Transporter.Transport.Free()
}