package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// RestartSimulation is the state of the action that restarts a simulation.
type RestartSimulation struct {
	platform platform.Platform
	services application.Services
	GroupID  simulations.GroupID
}

// NewRestartSimulation initializes a new state for restarting simulations.
func NewRestartSimulation(platform platform.Platform, services application.Services, groupID simulations.GroupID) *RestartSimulation {
	return &RestartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
