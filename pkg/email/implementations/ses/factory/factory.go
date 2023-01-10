package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/cloud/aws"
	"github.com/gazebo-web/cloudsim/v4/pkg/email/implementations/ses"
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
)

// NewFunc is the factory creation function for the SES Storage implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Initialize dependencies
	dependenciesInitFns := []func(config *Config, dependencies *Dependencies) error{
		initializeAPI,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	api := ses.NewEmailSender(
		typeDependencies.API,
		typeDependencies.Logger,
	)

	// Set output value
	if err := factory.SetValue(out, api); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}

// initializeAPI initializes the API dependency
func initializeAPI(config *Config, dependencies *Dependencies) error {
	if dependencies.API != nil {
		return nil
	}

	// Prepare config provider
	awsConfig := aws.Config{Region: config.Region}
	cp, err := aws.GetConfigProvider(awsConfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create API
	dependencies.API = ses.NewAPI(cp)

	return nil
}
