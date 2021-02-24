package nps

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// JobsStartSimulation groups the jobs needed to start a simulation.
var JobsStartSimulation = actions.Jobs{
	LaunchInstances,
	LaunchGazeboServerPod,
}

// JobsStopSimulation groups the jobs needed to stop a simulation.
var JobsStopSimulation = actions.Jobs{
}
