package simulator

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

type ISimulator interface {
	Launch()
	Stop()
	Restart()
}

// Config represents a set of options to configure a Simulator.
type Config struct {
	NamePrefix               string `env:"AWS_INSTANCE_NAME_PREFIX,required"`
	ShouldTerminateInstances bool   `env:"EC2_NODE_MGR_TERMINATE_INSTANCES" envDefault:"true"`
	IamInstanceProfile       string `env:"AWS_IAM_INSTANCE_PROFILE_ARN" envDefault:"arn:aws:iam::200670743174:instance-profile/cloudsim-ec2-node"`
	JoinCmd                  string `env:"KUBEADM_JOIN,required"`
	AvailableEC2Machines     int    `env:"IGN_EC2_MACHINES_LIMIT" envDefault:"-1"`
}

// Simulator is the responsible of creating the nodes and registering them in the kubernetes master.
type Simulator struct {
	orchestrator    *orchestrator.Kubernetes
	cloud           *cloud.AmazonWS
	config          Config
}

// New returns a new Simulator instance.
func New(kubernetes *orchestrator.Kubernetes, aws *cloud.AmazonWS) *Simulator {
	cfg := Config{}
	if err := env.Parse(cfg); err != nil {
		// TODO: Throw an error. Logger? Log Fatal?
	}
	m := Simulator{
		orchestrator: kubernetes,
		cloud:        aws,
		config:       cfg,
	}
	return &m
}

func (s *Simulator) Launch() {

}

func (s *Simulator) Stop() {

}

func (s *Simulator) Restart() {

}
