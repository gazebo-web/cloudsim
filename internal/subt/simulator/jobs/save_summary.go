package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// SaveSummary is a job in charge of persisting the summary of a certain simulation
var SaveSummary = &actions.Job{
	Name:       "save-simulation-summary",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    saveSummary,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// saveSummary is the main execute function for the SaveSummary job.
func saveSummary(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	summary, err := s.SubTServices().Summaries().Save(s.GroupID, s.Score, s.Stats, s.RunData)
	if err != nil {
		return nil, err
	}

	s.Summary = *summary
	store.SetState(s)

	return s, nil

}
