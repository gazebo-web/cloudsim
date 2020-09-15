package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// RestartSimulation is the state of the action that restarts a simulation.
type RestartSimulation struct {
	Platform platform.Platform
	Services application.Services
	GroupID  simulations.GroupID
}
