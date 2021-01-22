package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SetSimulationStatusToRunning is used to set a simulation status to running.
var SetSimulationStatusToRunning = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-running",
	Status:     simulations.StatusRunning,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
	PostHooks:  nil,
})

// SetSimulationStatusToLaunchInstances is used to set a simulation status to launch instances.
var SetSimulationStatusToLaunchInstances = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-launch-instances",
	Status:     simulations.StatusLaunchingInstances,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
	PostHooks:  nil,
})

// SetSimulationStatusToLaunchPods is used to set a simulation status to launching pods.
var SetSimulationStatusToLaunchPods = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-launch-pods",
	Status:     simulations.StatusLaunchingPods,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
	PostHooks:  nil,
})

// SetSimulationStatusToWaitInstances is used to set a simulation status to waiting instances.
var SetSimulationStatusToWaitInstances = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-wait-instances",
	Status:     simulations.StatusWaitingInstances,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
	PostHooks:  nil,
})

// SetSimulationStatusToWaitNodes is used to set a simulation status to waiting nodes.
var SetSimulationStatusToWaitNodes = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-wait-nodes",
	Status:     simulations.StatusWaitingNodes,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
	PostHooks:  nil,
})

// SetSimulationStatusToWaitPods is used to set a simulation status to waiting pods.
var SetSimulationStatusToWaitPods = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-wait-pods",
	Status:     simulations.StatusWaitingPods,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, returnGroupIDFromStartState},
})

// SetSimulationStatusToDeletingPods is used to set a simulation status to deleting pods.
var SetSimulationStatusToDeletingPods = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-deleting-pods",
	Status:     simulations.StatusDeletingPods,
	InputType:  &state.StopSimulation{},
	OutputType: &state.StopSimulation{},
	PreHooks:   []actions.JobFunc{setStopState, returnGroupIDFromStopState},
})

// SetSimulationStatusToDeletingNodes is used to set a simulation status to deleting nodes.
var SetSimulationStatusToDeletingNodes = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-deleting-nodes",
	Status:     simulations.StatusDeletingNodes,
	InputType:  &state.StopSimulation{},
	OutputType: &state.StopSimulation{},
	PreHooks:   []actions.JobFunc{setStopState, returnGroupIDFromStopState},
})

// SetSimulationStatusToProcessingResults is used to set a simulation status to processing results.
var SetSimulationStatusToProcessingResults = GenerateSetSimulationStatusJob(GenerateSetSimulationStatusConfig{
	Name:       "set-simulation-status-processing-results",
	Status:     simulations.StatusProcessingResults,
	InputType:  &state.StopSimulation{},
	OutputType: &state.StopSimulation{},
	PreHooks:   []actions.JobFunc{setStopState, returnGroupIDFromStopState},
})
