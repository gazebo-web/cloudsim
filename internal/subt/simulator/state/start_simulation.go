package state

import (
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// StartSimulation is the state of the action that starts a simulation.
type StartSimulation struct {
	state.PlatformGetter
	state.ServicesGetter
	platform             platform.Platform
	services             subtapp.Services
	GroupID              simulations.GroupID
	GazeboServerPod      orchestrator.Resource
	CreateMachinesInput  []cloud.CreateMachinesInput
	CreateMachinesOutput []cloud.CreateMachinesOutput
	ParentGroupID        *simulations.GroupID
	GazeboServerIP       string
}

// Platform returns the underlying platform.
func (s *StartSimulation) Platform() platform.Platform {
	return s.platform
}

// Services returns the underlying application services.
func (s *StartSimulation) Services() application.Services {
	return s.services
}

func (s *StartSimulation) SpecificServices() subtapp.Services {
	return s.services
}

// NewStartSimulation initializes a new state for starting simulations.
func NewStartSimulation(platform platform.Platform, services subtapp.Services, groupID simulations.GroupID) *StartSimulation {
	return &StartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
