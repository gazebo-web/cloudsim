package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// returnState is an actions.JobFunc implementation that returns the state. It's usually used as a posthook.
func returnState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	return store.State(), nil
}

// setStartState parses the input value as the StartSimulation state and sets the store with that state.
func setStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := value.(*state.StartSimulation)
	store.SetState(s)
	return s, nil
}

// returnGroupIDFromStartState parses the input value as the StartSimulation state and returns the group id.
func returnGroupIDFromStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(*state.StartSimulation)
	return s.GroupID, nil
}
