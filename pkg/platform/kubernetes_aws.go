package platform

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// InitializeAWS initializes the components from Amazon Web Services.
func InitializeAWS(region string, logger ign.Logger) (cloud.Storage, cloud.Machines, error) {
	config := aws.Config{Region: region}
	cp, err := aws.GetConfigProvider(config)
	if err != nil {
		return nil, nil, err
	}
	s3API := s3.NewAPI(cp)
	ec2API := ec2.NewAPI(cp)
	return s3.NewStorage(s3API, logger), ec2.NewMachines(ec2API, logger), nil
}

// InitializeKubernetes initializes a new Kubernetes orchestrator.
func InitializeKubernetes(logger ign.Logger) (orchestrator.Cluster, error) {
	config, err := kubernetes.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewAPI(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return kubernetes.NewDefaultKubernetes(client, spdyInit, logger), nil
}
