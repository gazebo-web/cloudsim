package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	jobs.CheckSimulationPendingStatus,
	jobs.CheckSimulationIsNotParent,
	jobs.CheckSimulationNoErrors,
	jobs.SetSimulationStatusToLaunchInstances,
	jobs.LaunchInstances,
	jobs.SetSimulationStatusToWaitInstances,
	jobs.WaitForInstances,
	jobs.SetSimulationStatusToWaitNodes,
	jobs.WaitForNodes,
	jobs.SetSimulationStatusToLaunchPods,
	jobs.CreateNetworkPolicyGazeboServer,
	jobs.LaunchGazeboServerPod,
	jobs.WaitForGazeboServerPod,
	jobs.WaitUpstream,
	jobs.ConfigureIngressGloo,
	jobs.LaunchWebsocketService,
	jobs.LaunchCommsBridgePods,
	jobs.CreateNetworkPolicyCommsBridges,
	jobs.LaunchFieldComputerPods,
	jobs.CreateNetworkPolicyFieldComputers,
	jobs.SetSimulationStatusToWaitPods,
	jobs.WaitRobots,
	jobs.SetSimulationStatusToRunning,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{}

// JobsRestartSimulation groups the jobs needed to restart a simulation.
var JobsRestartSimulation = actions.Jobs{}
