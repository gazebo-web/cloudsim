package jobs

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/actions"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/pkg/simulator"
	"github.com/gazebo-web/cloudsim/pkg/simulator/state"
	"github.com/jinzhu/gorm"
)

// CreateConfigurationsInput is the input of the CreateConfigurations job.
type CreateConfigurationsInput []configurations.CreateConfigurationInput

// CreateConfigurationsOutput is the output of the CreateConfigurations job.
// This struct was set in place to let the post-hook handle errors.
type CreateConfigurationsOutput struct {
	Resources []resource.Resource
	Error     error
}

const createdConfigurationsJobDataType = "created-configurations"

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
			return CreateConfigurationsOutput{}, nil
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
		res, err = s.Platform().Orchestrator().Configurations().Create(context.TODO(), in)
		if err != nil {
			return nil, err
		}
		created = append(created, res)
	}

	// Store configuration names
	configs := make([]string, 0, len(created))
	for _, res := range created {
		configs = append(configs, res.Name())
	}
	err = deployment.SetJobData(tx, nil, createdConfigurationsJobDataType, configs)

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

	// Get the set of configurations from the job data
	configs := make([]string, 0)
	dataErr := deployment.GetJobDataOutValue(tx, nil, createdConfigurationsJobDataType, &configs)
	if dataErr != nil {
		return nil, dataErr
	}

	// Delete the confingurations
	namespace := s.Platform().Store().Orchestrator().Namespace()
	for _, name := range configs {
		res := resource.NewResource(
			name,
			namespace,
			nil,
		)
		_, _ = s.Platform().Orchestrator().Configurations().Delete(context.TODO(), res)
	}

	return nil, nil
}
