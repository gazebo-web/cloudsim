package platform

import (
	"github.com/marcoshuck/cloudsim-refactor-proposal/pkg/platform/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type k8sAwsPlatform struct {
	storage      cloud.Storage
	machines     cloud.Machines
	orchestrator orchestrator.Orchestrator
}

func (p *k8sAwsPlatform) Storage() cloud.Storage {
	return p.storage
}

func (p *k8sAwsPlatform) Machines() cloud.Machines {
	return p.machines
}

func (p *k8sAwsPlatform) Orchestrator() orchestrator.Orchestrator {
	return p.orchestrator
}

// initializeAWS initializes the components from Amazon Web Services.
func initializeAWS(logger ign.Logger) (cloud.Storage, cloud.Machines, error) {
	// TODO: Read Region from env vars.
	config := aws.Config{Region: "us-east-1"}
	cp, err := aws.GetConfigProvider(config)
	if err != nil {
		return nil, nil, err
	}
	s3API := s3.GetClient(cp)
	ec2API := ec2.GetClient(cp)
	return s3.NewStorage(s3API, logger), ec2.NewMachines(ec2API, logger), nil
}

// initializeKubernetes initializes a new Kubernetes orchestrator.
func initializeKubernetes(logger ign.Logger) (orchestrator.Orchestrator, error) {
	config, err := kubernetes.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.GetClient(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return kubernetes.NewDefaultKubernetes(client, spdyInit, logger), nil
}

// NewAmazonWebServicesKubernetesPlatform initializes a new platform that uses AWS and Kubernetes.
func NewAmazonWebServicesKubernetesPlatform(logger ign.Logger) (Platform, error) {
	storage, machines, err := initializeAWS(logger)
	if err != nil {
		return nil, err
	}
	k8s, err := initializeKubernetes(logger)
	if err != nil {
		return nil, err
	}
	return &k8sAwsPlatform{
		storage:      storage,
		machines:     machines,
		orchestrator: k8s,
	}, nil
}
