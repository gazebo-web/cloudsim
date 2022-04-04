package jobs

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/state"
)

// CreateConfigurationsInput is the input of the CreateConfigurations job.
type CreateConfigurationsInput []configurations.CreateConfigurationInput

// CreateConfigurationsOutput is the output of the CreateConfigurations job.
// This struct was set in place to let the post-hook handle errors.
type CreateConfigurationsOutput struct {
	Resources []resource.Resource
	Error     error
}

// CreateConfigurations is a generic job to create cluster configurations.
var CreateConfigurations = &actions.Job{
	Name:       "create-configurations",
	Execute:    createConfigurations,
	InputType:  actions.GetJobDataType(&CreateConfigurationsInput{}),
	OutputType: actions.GetJobDataType(&CreateConfigurationsOutput{}),
}

// createConfigurations is the main function executed by the CreateConfigurations job.
func createConfigurations(store actions.Store, tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
	s := store.State().(state.PlatformGetter)

	// Parse input
	input, ok := value.(CreateConfigurationsInput)
	if !ok {
		// If assertion fails but CreateConfigurationsInput is nil, assume that no configurations need to be created.
		if input == nil {
			return CreateConfigurationsOutput{
				Resources: []resource.Resource{},
				Error:     nil,
			}, nil
		}

		return nil, simulator.ErrInvalidInput
	}

	if len(input) == 0 {
		return CreateConfigurationsOutput{
			Resources: []resource.Resource{},
			Error:     nil,
		}, nil
	}

	var created []resource.Resource
	var err error

	for _, in := range input {
		var res resource.Resource
		res, err = s.Platform().Orchestrator().Configurations().Create(in)
		if err != nil {
			return nil, err
		}
		created = append(created, res)
	}

	return CreateConfigurationsOutput{
		Resources: created,
		Error:     err,
	}, nil
}

// DeleteCreatedConfigurationsOnFailure is an optional rollback handler that removes any created configurations when
// an action fails.
func DeleteCreatedConfigurationsOnFailure(store actions.Store, tx *gorm.DB, deployment *actions.Deployment,
	value interface{}, err error) (interface{}, error) {

	// Get the store
	s := store.State().(state.PlatformGetter)

	// Get the list of instances from the execute function
	data, dataErr := deployment.GetJobData(tx, nil, actions.DeploymentJobInput)
	if dataErr != nil {
		return nil, dataErr
	}

	// Parse the list of instances
	createdConfigurations, ok := data.(*CreateConfigurationsInput)
	if !ok || createdConfigurations == nil {
		return nil, err
	}

	// Terminate the instances
	for _, c := range *createdConfigurations {
		res := resource.NewResource(c.Name, c.Namespace, nil)
		_, _ = s.Platform().Orchestrator().Configurations().Delete(res)
	}

	return nil, err
}
