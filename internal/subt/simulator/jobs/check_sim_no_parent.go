package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckStartSimulationIsNotParent is a job in charge of checking if a simulation is not a parent simulation.
var CheckStartSimulationIsNotParent = GenerateCheckSimulationNotOfKindJob(
	ConfigCheckSimulationNotOfKindJob{
		Name:               "check-start-simulation-no-parent",
		Kind:               simulations.SimParent,
		PreHooks:           []actions.JobFunc{setStartState},
		PreparationPreHook: generateCheckStartSimulationNotOfKindInputPreHook(simulations.SimParent),
		InputType:          &state.StartSimulation{},
		OutputType:         &state.StartSimulation{},
	},
)

// CheckStopSimulationIsNotParent is a job in charge of checking if a simulation is not a parent simulation.
var CheckStopSimulationIsNotParent = GenerateCheckSimulationNotOfKindJob(
	ConfigCheckSimulationNotOfKindJob{
		Name:               "check-stop-simulation-no-parent",
		Kind:               simulations.SimParent,
		PreHooks:           []actions.JobFunc{setStopState},
		PreparationPreHook: generateCheckStopSimulationNotOfKindInputPreHook(simulations.SimParent),
		InputType:          &state.StopSimulation{},
		OutputType:         &state.StopSimulation{},
	},
)
