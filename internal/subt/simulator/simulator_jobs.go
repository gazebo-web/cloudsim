package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	jobs.CheckSimulationPendingStatus,
	jobs.CheckStartSimulationIsNotParent,
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
	jobs.GetGazeboIP,
	jobs.WaitUpstream,
	jobs.ConfigureIngressGloo,
	jobs.LaunchWebsocketService,
	jobs.LaunchCommsBridgePods,
	jobs.CreateNetworkPolicyCommsBridges,
	jobs.LaunchFieldComputerPods,
	jobs.CreateNetworkPolicyFieldComputers,
	jobs.SetSimulationStatusToWaitPods,
	jobs.WaitSimulationPods,
	jobs.AddRunningSimulation,
	jobs.SetSimulationStatusToRunning,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{
	jobs.CheckSimulationTerminateRequestedStatus,
	jobs.SetSimulationStatusToProcessingResults,
	jobs.CheckStopSimulationIsNotParent,
	jobs.SetStoppedAt,
	jobs.ReadScore,
	jobs.ReadStats,
	jobs.ReadRunData,
	// jobs.GenerateSummary,
}
