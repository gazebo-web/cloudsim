package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StartSimulation is the state of the action that starts a simulation.
type StartSimulation struct {
	Platform platform.Platform
	Services application.Services
	GroupID  simulations.GroupID
}
