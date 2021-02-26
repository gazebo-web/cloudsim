package nps

import (
  "fmt"
  "github.com/jinzhu/gorm"
  "gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// returnState is an actions.JobFunc implementation that returns the state. It's usually used as a posthook.
func returnState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\nreturnState\n")
  return store.State(), nil
}

// setStartState parses the input value as the StartSimulation state and sets the store with that state.
func setStartState(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\nsetStartState\n")
	s := value.(*StartSimulationData)
	store.SetState(s)
	return s, nil
}

// checkWaitError is an actions.JobFunc implementation that checks if the value returned by a job extended from the
// Wait job returns an error.
// If the previous jobs.Wait job returns an error, it will trigger the rollback handler.
// If the previous jobs.Wait does not return an error, it will pass to the next actions.JobFunc.
func checkWaitError(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
  fmt.Printf("\ncheckWaitError\n")
	output := value.(jobs.WaitOutput)
	if output.Error != nil {
		return nil, output.Error
	}
	return nil, nil
}
