package jobs

import (
	"github.com/jinzhu/gorm"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/runsim"
)

// AddRunningSimulation is job in charge of adding a running simulation to the list of running simulations.
var AddRunningSimulation = &actions.Job{
	Name:            "add-running-simulation",
	PreHooks:        []actions.JobFunc{setStartState},
	Execute:         addRunningSimulation,
	PostHooks:       []actions.JobFunc{returnState},
	RollbackHandler: revertAddingRunningSimulation,
	InputType:       actions.GetJobDataType(&state.StartSimulation{}),
	OutputType:      actions.GetJobDataType(&state.StartSimulation{}),
}

// revertAddingRunningSimulation reverts all the changes made while adding a running simulation.
func revertAddingRunningSimulation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}, _ error) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	s.Platform().RunningSimulations().Free(s.GroupID)

	err := s.Platform().RunningSimulations().Remove(s.GroupID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// addRunningSimulation is the main function of the AddRunningSimulation job.
func addRunningSimulation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)

	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	t, err := s.SubTServices().Tracks().Get(subtSim.GetTrack())
	if err != nil {
		return nil, err
	}

	rs := runsim.NewRunningSimulation(s.GroupID, int64(t.MaxSimSeconds), sim.GetValidFor())

	err = s.Platform().RunningSimulations().Add(s.GroupID, rs, s.WebsocketConnection)
	if err != nil {
		return nil, err
	}

	return s, nil
}
