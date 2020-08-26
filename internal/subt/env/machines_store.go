package env

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

type machineEnvStore struct {
	InstanceProfileValue string   `env:"SUBT_MACHINES_INSTANCE_PROFILE" envDefault:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker"`
	KeyNameValue         string   `env:"SUBT_MACHINES_KEY_NAME" envDefault:"ignitionFuel"`
	MachineTypeValue     string   `env:"SUBT_MACHINES_TYPE" envDefault:"g3.4xlarge"`
	FirewallRulesValue   []string `env:"SUBT_MACHINES_FIREWALL_RULES" envSeparator:"," envDefault:"sg-0c5c791266694a3ca"`
	SubnetsValue         []string `env:"SUBT_MACHINES_SUBNETS,required" envSeparator:","`
	ZonesValue           []string `env:"SUBT_MACHINES_ZONES,required" envSeparator:","`
	MachinesLimitValue   int      `env:"SUBT_MACHINES_LIMIT" envDefault:"-1"`
	createdMachines      int
}

func (m *machineEnvStore) InstanceProfile() *string {
	return &m.InstanceProfileValue
}

func (m *machineEnvStore) KeyName() string {
	return m.KeyNameValue
}

func (m *machineEnvStore) Type() string {
	return m.MachineTypeValue
}

func (m *machineEnvStore) FirewallRules() []string {
	return m.FirewallRulesValue
}

func (m *machineEnvStore) Subnet() string {
	i := m.createdMachines % len(m.SubnetsValue)
	return m.SubnetsValue[i]
}

func (m *machineEnvStore) Zone() string {
	i := m.createdMachines % len(m.ZonesValue)
	return m.SubnetsValue[i]
}

func (m *machineEnvStore) Tags(simulation simulations.Simulation) []cloud.Tag {
	return nil
}

func (m *machineEnvStore) InitScript() *string {
	return nil
}

func newMachinesStore() store.Machines {
	var m machineEnvStore
	if err := env.Parse(&m); err != nil {
		panic(err)
	}
	return &m
}
