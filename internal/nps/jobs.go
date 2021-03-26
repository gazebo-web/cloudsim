package nps

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// StartSimulationAction groups the jobs needed to start a simulation.
// See also StartSimulationData, which holds information passed between jobs.
var StartSimulationAction = actions.Jobs{
	LaunchInstances,
	WaitForInstances,
	WaitForNodes,
	LaunchPod,
	WaitForPod,
	GetPodIP,
}

// StartSimulationData contains information is that is passed to each job in
// the StartSimulationAction.
type StartSimulationData struct {
	state.PlatformGetter
	state.ServicesGetter
	platform platform.Platform
	GroupID  simulations.GroupID
	URI      string
	// NodeSelector allows a job to select the correct running kubernetes  node.
	NodeSelector orchestrator.Selector
	// PodSelector allows a job to select the correct running kubernetes pod.
	PodSelector orchestrator.Selector
	logger      ign.Logger
	// \todo: What is this used for? I'm using it launch_instance_job.go for some reason.
	CreateMachinesInput []cloud.CreateMachinesInput
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
	RemovePods,
	RemoveInstances,
}

// StopSimulationData contains information is that is passed to each job in
// the StopSimulationAction.
type StopSimulationData struct {
	state.PlatformGetter
	state.ServicesGetter
	platform platform.Platform
	GroupID  simulations.GroupID
	// NodeSelector allows a job to select the correct running kubernetes  node.
	NodeSelector orchestrator.Selector
	// PodSelector allows a job to select the correct running kubernetes pod.
	PodSelector orchestrator.Selector
	logger      ign.Logger
	PodList     []orchestrator.Resource
}

// Platform returns the underlying platform contained in StartSimulationData
// See state.PlatformGetter in the StartSimulationData struct.
//
// \todo Major Improvement needed: If this function is not defined, then a job can fail with little to no meaningful output. The only way to know it's missing is to dive into the job's implementation, see that the job fails when calling Platform() and remember that you forgot to implement this function. This pattern needd to be improved.
func (s *StopSimulationData) Platform() platform.Platform {
	return s.platform
}
