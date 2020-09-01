package aws

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// InitializeAWS initializes all the components used from Amazon Web Services.
func InitializeAWS(region string, logger ign.Logger) (cloud.Storage, cloud.Machines, error) {
	config := Config{Region: region}
	cp, err := GetConfigProvider(config)
	if err != nil {
		return nil, nil, err
	}
	s3API := s3.NewAPI(cp)
	ec2API := ec2.NewAPI(cp)
	return s3.NewStorage(s3API, logger), ec2.NewMachines(ec2API, logger), nil
}
