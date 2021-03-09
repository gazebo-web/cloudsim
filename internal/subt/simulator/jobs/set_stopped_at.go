package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// SetStoppedAt is a job in charge of setting StoppedAt field from a certain simulation to the time of when it has been stopped.
var SetStoppedAt = actions.Job{
	Name:       "set-stopped-at",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    setStoppedAt,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// setStoppedAt is the execute function of the SetStoppedAt job that will set the simulation's StoppedAt value.
func setStoppedAt(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	err := s.SubTServices().Simulations().MarkStopped(s.GroupID)

	if err != nil {
		return nil, err
	}

	return s, nil
}
