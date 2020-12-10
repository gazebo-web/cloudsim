package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	jobs.CheckSimulationPendingStatus,
	jobs.CheckSimulationIsParent,
	jobs.CheckSimulationNoErrors,
	jobs.SetSimulationStatusToLaunchInstances,
	jobs.LaunchInstances,
	jobs.SetSimulationStatusToWaitInstances,
	jobs.WaitForInstances,
	jobs.SetSimulationStatusToWaitNodes,
	jobs.WaitForNodes,
	jobs.SetSimulationStatusToLaunchPods,
	// jobs.CreateGazeboServerNetworkPolicy,
	jobs.LaunchGazeboServerPod,
	jobs.WaitForGazeboServerPod,
	// jobs.LaunchGazeboServerCopyPod,
	jobs.WaitUpstream,
	jobs.ConfigureIngressGloo,
	// jobs.LaunchWebsocketService,
	// jobs.CreateRobotNetworkPolicies,
	// jobs.LaunchCommsBridgePods,
	// jobs.LaunchFieldComputerPods,
	jobs.SetSimulationStatusToWaitPods,
	//jobs.WaitForRobotPods,
	jobs.SetSimulationStatusToRunning,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{}

// JobsRestartSimulation groups the jobs needed to restart a simulation.
var JobsRestartSimulation = actions.Jobs{}
