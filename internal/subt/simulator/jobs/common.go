package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
)

// returnState is an actions.JobFunc implementation that returns the state. It's usually used as a posthook.
func returnState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	return store.State(), nil
}

// checkError is an actions.JobFunc implementation that checks if the value given by the previous actions.JobFunc returns an error.
// If the previous actions.JobFunc returns an error, it will trigger the rollback handler.
// If the previous actions.JobFunc does not return an error, it will pass to the next actions.JobFunc.
func checkError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	err := value.(error)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
