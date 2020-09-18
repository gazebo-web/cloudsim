package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// GenerateSetSimulationStatusJob generates a job to set a simulation to a certain status.
func GenerateSetSimulationStatusJob(name string, status simulations.Status, inputType, outputType interface{}, prehooks ...actions.JobFunc) *actions.Job {
	return jobs.SetSimulationStatus.Extend(actions.Job{
		Name:       name,
		PreHooks:   append(prehooks, generateSetSimulationStatusInputPreHook(status)),
		PostHooks:  []actions.JobFunc{returnState},
		InputType:  actions.GetJobDataType(inputType),
		OutputType: actions.GetJobDataType(outputType),
	})
}

func generateSetSimulationStatusInputPreHook(status simulations.Status) actions.JobFunc {
	return func(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		gid := value.(simulations.GroupID)

		return jobs.SetSimulationStatusInput{
			GroupID: gid,
			Status:  status,
		}, nil
	}
}
