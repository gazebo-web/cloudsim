package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StopSimulation is the state of the action that stops a simulation.
type StopSimulation struct {
	platform platform.Platform
	services application.Services
	GroupID  simulations.GroupID
}

// NewStopSimulation initializes a new state for stopping simulations.
func NewStopSimulation(platform platform.Platform, services application.Services, groupID simulations.GroupID) *RestartSimulation {
	return &RestartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
