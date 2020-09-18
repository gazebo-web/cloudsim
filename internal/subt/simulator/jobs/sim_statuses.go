package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SetSimulationStatusToRunning is used to set a simulation status to running.
var SetSimulationStatusToRunning = GenerateSetSimulationStatusJob(
	"set-simulation-status-running",
	simulations.StatusRunning,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)

// SetSimulationStatusToLaunchInstances is used to set a simulation status to launch instances.
var SetSimulationStatusToLaunchInstances = GenerateSetSimulationStatusJob(
	"set-simulation-status-launch-instances",
	simulations.StatusLaunchingInstances,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)

// SetSimulationStatusToLaunchPods is used to set a simulation status to launching pods.
var SetSimulationStatusToLaunchPods = GenerateSetSimulationStatusJob(
	"set-simulation-status-launch-pods",
	simulations.StatusLaunchingPods,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)

// SetSimulationStatusToWaitInstances is used to set a simulation status to waiting instances.
var SetSimulationStatusToWaitInstances = GenerateSetSimulationStatusJob(
	"set-simulation-status-wait-instances",
	simulations.StatusWaitingInstances,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)

// SetSimulationStatusToWaitNodes is used to set a simulation status to waiting nodes.
var SetSimulationStatusToWaitNodes = GenerateSetSimulationStatusJob(
	"set-simulation-status-wait-nodes",
	simulations.StatusWaitingNodes,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)

// SetSimulationStatusToWaitPods is used to set a simulation status to waiting pods.
var SetSimulationStatusToWaitPods = GenerateSetSimulationStatusJob(
	"set-simulation-status-wait-pods",
	simulations.StatusWaitingPods,
	&state.StartSimulation{},
	&state.StartSimulation{},
	setStartState,
	returnGroupIDFromStartState,
)
