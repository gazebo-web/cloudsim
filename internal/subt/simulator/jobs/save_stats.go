package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// SaveStats is a job in charge of persisting the stats from a certain simulation
var SaveStats = actions.Job{
	Name:       "save-simulation-stats",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    saveStats,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// saveStats is the main execute function for the SaveStats job.
func saveStats(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	err := s.SubTServices().Statistics().Save(s.GroupID, &s.Score, s.Stats)
	if err != nil {
		return nil, err
	}

	return s, nil
}
