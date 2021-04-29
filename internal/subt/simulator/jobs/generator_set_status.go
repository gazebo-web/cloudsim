package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/jobs"
)

// GenerateSetSimulationStatusConfig is the config needed by GenerateSetSimulationStatusJob to configure a new job.
type GenerateSetSimulationStatusConfig struct {
	Name       string
	Status     simulations.Status
	InputType  interface{}
	OutputType interface{}
	PreHooks   []actions.JobFunc
	PostHooks  []actions.JobFunc
}

// GenerateSetSimulationStatusJob generates a job from jobs.SetSimulationStatus to set a certain status to a simulation.
func GenerateSetSimulationStatusJob(config GenerateSetSimulationStatusConfig) *actions.Job {
	return jobs.SetSimulationStatus.Extend(actions.Job{
		Name:       config.Name,
		PreHooks:   append(config.PreHooks, generateSetSimulationStatusInputPreHook(config.Status)),
		PostHooks:  append(config.PostHooks, returnState),
		InputType:  actions.GetJobDataType(config.InputType),
		OutputType: actions.GetJobDataType(config.OutputType),
	})
}

// generateSetSimulationStatusInputPreHook is a pre-hook in charge of preparing the input for the generic
// jobs.SetSimulationStatus job
func generateSetSimulationStatusInputPreHook(status simulations.Status) actions.JobFunc {
	return func(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		gid := value.(simulations.GroupID)

		return jobs.SetSimulationStatusInput{
			GroupID: gid,
			Status:  status,
		}, nil
	}
}
