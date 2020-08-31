package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/jobs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	jobs.CheckPendingStatus,
	jobs.CheckSimulationParenthood,
	jobs.CheckParentSimulationWithError,
	jobs.UpdateSimulationStatusToLaunchInstances,
	jobs.LaunchInstances,
	jobs.UpdateSimulationStatusToWaitInstances,
	jobs.WaitForInstances,
	jobs.UpdateSimulationStatusToInstancesReady,
	jobs.UpdateSimulationStatusToWaitNodes,
	jobs.WaitForKubernetesNodes,
	jobs.UpdateSimulationStatusToNodesReady,
	jobs.UpdateSimulationStatusToLaunchPods,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{}

// JobsRestartSimulation groups the jobs needed to restart a simulation.
var JobsRestartSimulation = actions.Jobs{}
