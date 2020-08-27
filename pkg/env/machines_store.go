package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// machineEnvStore is a store.Machines implementation.
// It contains all the information needed by the SubT jobs related to machines.
type machineEnvStore struct {
	InstanceProfileValue string   `env:"CLOUDISM_MACHINES_INSTANCE_PROFILE" envDefault:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker"`
	KeyNameValue         string   `env:"CLOUDISM_MACHINES_KEY_NAME" envDefault:"ignitionFuel"`
	MachineTypeValue     string   `env:"CLOUDISM_MACHINES_TYPE" envDefault:"g3.4xlarge"`
	FirewallRulesValue   []string `env:"CLOUDISM_MACHINES_FIREWALL_RULES" envSeparator:"," envDefault:"sg-0c5c791266694a3ca"`
	SubnetsValue         []string `env:"CLOUDISM_MACHINES_SUBNETS,required" envSeparator:","`
	ZonesValue           []string `env:"CLOUDISM_MACHINES_ZONES,required" envSeparator:","`
	MachinesLimitValue   int      `env:"CLOUDISM_MACHINES_LIMIT" envDefault:"-1"`
	BaseImageValue       string   `env:"CLOUDISM_MACHINES_BASE_IMAGE" envDefault:"ami-08861f7e7b409ed0c"`
	NamePrefixValue      string   `env:"CLOUDISM_MACHINES_NAME_PREFIX,required" envDefault:"cloudsim-subt-node"`
	ClusterNameValue     string   `env:"CLOUDISM_MACHINES_CLUSTER_NAME,required"`
	NodeReadyTimeout     uint     `env:"CLOUDSIM_MACHINES_NODE_READY_TIMEOUT_SECONDS" envDefault:"300"`
	subnetZoneIndex      int
}

// BaseImage returns the base image value.
// In SubT, the base image is the Amazon
func (m *machineEnvStore) BaseImage() string {
	return m.BaseImageValue
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

func (m *machineEnvStore) subnet() string {
	i := m.subnetZoneIndex % len(m.SubnetsValue)
	return m.SubnetsValue[i]
}

func (m *machineEnvStore) zone() string {
	i := m.subnetZoneIndex % len(m.ZonesValue)
	return m.SubnetsValue[i]
}

func (m *machineEnvStore) Timeout() time.Duration {
	return time.Duration(m.NodeReadyTimeout) * time.Second
}

func (m *machineEnvStore) PollFrequency() time.Duration {
	return 2 * time.Second
}

func (m *machineEnvStore) SubnetAndZone() (string, string) {
	subnet, zone := m.subnet(), m.zone()
	m.subnetZoneIndex++
	return subnet, zone
}

func (m *machineEnvStore) Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []cloud.Tag {
	name := fmt.Sprintf("%s-%s-%s", m.NamePrefixValue, simulation.GroupID(), nameSuffix)
	clusterKey := fmt.Sprintf("kubernetes.io/cluster/%s", m.ClusterNameValue)
	return []cloud.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                       name,
				"cloudsim_groupid":           string(simulation.GroupID()),
				"CloudsimGroupID":            string(simulation.GroupID()),
				"project":                    "cloudsim",
				"Cloudsim":                   "True",
				"SubT":                       "True",
				"cloudsim-application":       "SubT",
				"cloudsim-simulation-worker": m.NamePrefixValue,
				"cloudsim_node_type":         nodeType,
				clusterKey:                   "owned",
			},
		},
	}
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
