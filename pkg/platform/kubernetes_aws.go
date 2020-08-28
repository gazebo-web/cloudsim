package platform

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// k8sAwsPlatform is a platform implementation using AWS and Kubernetes.
type k8sAwsPlatform struct {
	storage      cloud.Storage
	machines     cloud.Machines
	orchestrator orchestrator.Cluster
	store        store.Store
}

// Store returns a store.Store implementation.
func (p *k8sAwsPlatform) Store() store.Store {
	return p.store
}

// Storage returns a cloud.Storage implementation.
func (p *k8sAwsPlatform) Storage() cloud.Storage {
	return p.storage
}

// Machines returns a cloud.Machines implementation.
func (p *k8sAwsPlatform) Machines() cloud.Machines {
	return p.machines
}

// Orchestrator returns an orchestrator.Cluster implementation.
func (p *k8sAwsPlatform) Orchestrator() orchestrator.Cluster {
	return p.orchestrator
}

// InitializeAWS initializes the components from Amazon Web Services.
func InitializeAWS(region string, logger ign.Logger) (cloud.Storage, cloud.Machines, error) {
	config := aws.Config{Region: region}
	cp, err := aws.GetConfigProvider(config)
	if err != nil {
		return nil, nil, err
	}
	s3API := s3.GetClient(cp)
	ec2API := ec2.GetClient(cp)
	return s3.NewStorage(s3API, logger), ec2.NewMachines(ec2API, logger), nil
}

// InitializeKubernetes initializes a new Kubernetes orchestrator.
func InitializeKubernetes(logger ign.Logger) (orchestrator.Cluster, error) {
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
func NewAmazonWebServicesKubernetesPlatform(ec2 cloud.Machines, s3 cloud.Storage, k8s orchestrator.Cluster, store store.Store) Platform {
	return &k8sAwsPlatform{
		storage:      s3,
		machines:     ec2,
		orchestrator: k8s,
		store:        store,
	}
}
