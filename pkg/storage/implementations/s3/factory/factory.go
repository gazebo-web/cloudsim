package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations/s3"
)

// NewFunc is the factory creation function for the S3 Storage implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return err
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
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
	api := s3.NewStorage(
		typeDependencies.API,
		typeDependencies.Logger,
	)

	// Set output value
	if err := factory.SetValue(out, api); err != nil {
		return err
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
		return err
	}

	// Create API
	dependencies.API = s3.NewAPI(cp)

	return nil
}
