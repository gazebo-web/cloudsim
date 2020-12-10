package state

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// StartSimulation is the state of the action that starts a simulation.
type StartSimulation struct {
	platform             platform.Platform
	services             application.Services
	GroupID              simulations.GroupID
	GazeboServerPod      orchestrator.Resource
	CreateMachinesInput  []cloud.CreateMachinesInput
	CreateMachinesOutput []cloud.CreateMachinesOutput
	ParentGroupID        *simulations.GroupID
	UpstreamName         string
}

// Platform returns the underlying platform.
func (s *StartSimulation) Platform() platform.Platform {
	return s.platform
}

// Services returns the underlying application services.
func (s *StartSimulation) Services() application.Services {
	return s.services
}

// NewStartSimulation initializes a new state for starting simulations.
func NewStartSimulation(platform platform.Platform, services application.Services, groupID simulations.GroupID) *StartSimulation {
	return &StartSimulation{
		platform: platform,
		services: services,
		GroupID:  groupID,
	}
}
