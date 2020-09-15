package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StopSimulation is the state of the action that stops a simulation.
type StopSimulation struct {
	Platform platform.Platform
	Services application.Services
	GroupID  simulations.GroupID
}
