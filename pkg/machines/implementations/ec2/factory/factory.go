package factory

import (
	"github.com/aws/aws-sdk-go/service/pricing"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations/ec2"
)

// NewFunc is the factory creation function for the EC2 Machines implementation.
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
		initializePricingAPI,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	api, err := ec2.NewMachines(&ec2.NewInput{
		API:             typeDependencies.API,
		Logger:          typeDependencies.Logger,
		Limit:           typeConfig.Limit,
		WorkerGroupName: typeConfig.WorkerGroupName,
		Region:          typeConfig.Region,
		Zones:           typeConfig.Zones,
		CostCalculator:  aws.NewCostCalculator(typeDependencies.PricingAPI, aws.ParseEC2, aws.KindMachines),
	})
	if err != nil {
		return err
	}

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
	dependencies.API = ec2.NewAPI(cp)

	return nil
}

// initializePricingAPI initializes the Pricing API dependency.
func initializePricingAPI(config *Config, dependencies *Dependencies) error {
	if dependencies.PricingAPI != nil {
		return nil
	}

	// Prepare config provider
	awsConfig := aws.Config{Region: config.Region}
	cp, err := aws.GetConfigProvider(awsConfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create API
	dependencies.PricingAPI = pricing.New(cp)

	return nil
}
