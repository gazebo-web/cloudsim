package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// CheckSimulationPendingStatus is used to check that a certain simulation has the pending status.
var CheckSimulationPendingStatus = GenerateCheckStatusJob(CheckStatusJobConfig{
	Name:       "check-simulation-pending-status",
	Status:     simulations.StatusPending,
	InputType:  &state.StartSimulation{},
	OutputType: &state.StartSimulation{},
	PreHooks:   []actions.JobFunc{setStartState, generateCheckStartSimulationStatusInputPreHook(simulations.StatusPending)},
})

// CheckSimulationTerminateRequestedStatus is used to check that a certain simulation has the terminate requested status.
var CheckSimulationTerminateRequestedStatus = GenerateCheckStatusJob(CheckStatusJobConfig{
	Name:       "check-simulation-terminate-requested-status",
	Status:     simulations.StatusTerminateRequested,
	InputType:  &state.StopSimulation{},
	OutputType: &state.StopSimulation{},
	PreHooks:   []actions.JobFunc{setStopState, generateCheckStopSimulationStatusInputPreHook(simulations.StatusTerminateRequested)},
})
