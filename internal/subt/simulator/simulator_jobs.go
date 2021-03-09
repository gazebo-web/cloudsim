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
	jobs.LaunchWebsocketService,
	jobs.WaitUpstream,
	jobs.ConfigureIngressGloo,
	jobs.CreateNetworkPolicyCommsBridges,
	jobs.CreateNetworkPolicyFieldComputers,
	jobs.LaunchCommsBridgePods,
	jobs.LaunchCommsBridgeCopyPods,
	jobs.WaitForCommsBridgePods,
	jobs.GetCommsBridgePodIP,
	jobs.LaunchFieldComputerPods,
	jobs.SetSimulationStatusToWaitPods,
	jobs.WaitSimulationPods,
	jobs.SetWebsocketConnection,
	jobs.AddRunningSimulation,
	jobs.SetSimulationStatusToRunning,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{
	jobs.CheckSimulationTerminateRequestedStatus,
	jobs.CheckSimulationNoErrors,
}
