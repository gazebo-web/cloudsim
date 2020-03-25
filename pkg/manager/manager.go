package manager

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

type Config struct {
	NamePrefix string `env:"AWS_INSTANCE_NAME_PREFIX,required"`
	ShouldTerminateInstances bool   `env:"EC2_NODE_MGR_TERMINATE_INSTANCES" envDefault:"true"`
	IamInstanceProfile       string `env:"AWS_IAM_INSTANCE_PROFILE_ARN" envDefault:"arn:aws:iam::200670743174:instance-profile/cloudsim-ec2-node"`
	JoinCmd string `env:"KUBEADM_JOIN,required"`
	AvailableEC2Machines int `env:"IGN_EC2_MACHINES_LIMIT" envDefault:"-1"`
}

type Manager struct {
	kc  *orchestrator.Kubernetes
	aws  *cloud.AmazonWS
	config Config
}

func New(kubernetes *orchestrator.Kubernetes, aws *cloud.AmazonWS) *Manager {
	cfg := Config{}
	if err := env.Parse(cfg); err != nil {
		// TODO: Throw an error
	}
	m := Manager{
		kc:     kubernetes,
		aws:    aws,
		config: cfg,
	}
	return &m
}