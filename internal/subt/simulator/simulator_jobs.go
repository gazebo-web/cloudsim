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
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{
	jobs.CheckSimulationTerminateRequestedStatus,
	jobs.CheckSimulationNoErrors,
}
