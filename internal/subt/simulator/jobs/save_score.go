package jobs

import (
	"github.com/jinzhu/gorm"
	subt "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SaveScore is a job in charge of persisting the score from a certain simulation
var SaveScore = &actions.Job{
	Name:       "save-simulation-score",
	PreHooks:   []actions.JobFunc{setStopState},
	Execute:    saveScore,
	PostHooks:  []actions.JobFunc{returnState},
	InputType:  actions.GetJobDataType(&state.StopSimulation{}),
	OutputType: actions.GetJobDataType(&state.StopSimulation{}),
}

// saveScore is the main execute function for the SaveScore job.
func saveScore(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StopSimulation)

	err := s.SubTServices().Simulations().UpdateScore(s.GroupID, s.Score)
	if err != nil {
		return nil, err
	}

	sim, err := s.SubTServices().Simulations().Get(s.GroupID)
	if err != nil {
		return nil, err
	}

	subtSim := sim.(subt.Simulation)

	applicationName := "subt"
	track := subtSim.GetTrack()

	gid := s.GroupID.String()

	if sim.IsKind(simulations.SimChild) {
		parent, err := s.SubTServices().Simulations().GetParent(s.GroupID)
		if err != nil {
			return nil, err
		}

		parentGroupID := parent.GetGroupID().String()

		em := s.SubTServices().Users().AddScore(&parentGroupID, &applicationName, &track, sim.GetOwner(), s.Score, &gid)
		if em != nil {
			return nil, em.BaseError
		}
		return s, nil
	}

	em := s.SubTServices().Users().AddScore(&gid, &applicationName, &track, sim.GetOwner(), s.Score, &gid)
	if em != nil {
		return nil, em.BaseError
	}

	return s, nil

}
