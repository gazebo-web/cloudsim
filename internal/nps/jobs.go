package nps

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
)

// StartSimulationAction groups the jobs needed to start a simulation.
// See also StartSimulationData, which holds information passed between jobs.
var StartSimulationAction = actions.Jobs{
	LaunchInstances,
  WaitForInstances,
  WaitForNodes,
	LaunchGazeboServerPod,
  WaitForPod,
  GetPodIP,
}

// StartSimulationData contains information is that is passed to each job in
// the StartSimulationAction.
type StartSimulationData struct {
	state.PlatformGetter
	state.ServicesGetter
  platform             platform.Platform
  GroupID              simulations.GroupID
  IP                   string

  // \todo: What is this used for? I'm using it launch_instance_job.go for some reason.
  CreateMachinesInput  []cloud.CreateMachinesInput
  // \todo: What is this used for? I'm using it launch_instance_job.go for some reason.
  CreateMachinesOutput []cloud.CreateMachinesOutput
}

// Platform returns the underlying platform contained in StartSimulationData
// See state.PlatformGetter in the StartSimulationData struct.
func (s *StartSimulationData) Platform() platform.Platform {
  return s.platform
}

// StopSimulationAction groups the jobs needed to stop a simulation.
var StopSimulationAction = actions.Jobs{
}

// StopSimulationData contains information is that is passed to each job in
// the StopSimulationAction.
type StopSimulationData struct {
}
