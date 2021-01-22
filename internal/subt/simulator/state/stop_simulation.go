package state

import (
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StopSimulation is the state of the action that stops a simulation.
type StopSimulation struct {
	platform platform.Platform
	services subtapp.Services
	GroupID  simulations.GroupID
	Score    float64
	Stats    simulations.Statistics
	RunData  string
}

// NewStopSimulation initializes a new state for stopping simulations.
func NewStopSimulation(platform platform.Platform, services subtapp.Services, groupID simulations.GroupID) *RestartSimulation {
	return &RestartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
