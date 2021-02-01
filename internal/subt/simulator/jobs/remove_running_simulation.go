package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// RemoveRunningSimulation is job in charge of removing a running simulation from the list of running simulations.
var RemoveRunningSimulation = &actions.Job{
	Name:       "add-running-simulation",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    removeRunningSimulation,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// removeRunningSimulation is the execute function for the RemoveRunningSimulation job.
func removeRunningSimulation(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	err := s.Platform().RunningSimulations().Remove(s.GroupID)
	if err != nil {
		return nil, err
	}

	return s, nil
}
