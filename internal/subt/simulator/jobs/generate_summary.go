package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// GenerateSummary is a job in charge of generating the summary for a certain simulation.
var GenerateSummary = &actions.Job{
	Name:       "generate-summary",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    generateSummary,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// generateSummary is the main execute function of the GenerateSummary job. It's in charge of getting all processed data in previous
// jobs and save a data structure with the summary.
func generateSummary(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	s.Summary = simulations.Summary{
		Score:                  s.Score,
		SimTimeDurationAvg:     float64(s.Stats.SimulationTime),
		SimTimeDurationStdDev:  0,
		RealTimeDurationAvg:    float64(s.Stats.RealTime),
		RealTimeDurationStdDev: 0,
		ModelCountAvg:          float64(s.Stats.ModelCount),
		ModelCountStdDev:       0,
		Sources:                "",
	}

	store.SetState(s)

	return store, nil
}
