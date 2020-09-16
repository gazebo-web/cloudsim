package jobs

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// getStartDataFromJob is used to get the current job data from store.
func getStartDataFromJob(store actions.Store) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	return s, nil
}
