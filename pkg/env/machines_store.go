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

// BaseImage returns the base image value read from env vars.
// In SubT, the base image is the Amazon
func (m *machineEnvStore) BaseImage() string {
	return m.BaseImageValue
}

// InstanceProfile returns the instance profile value read from env vars.
func (m *machineEnvStore) InstanceProfile() *string {
	return &m.InstanceProfileValue
}

// KeyName returns the key name value read from env vars.
func (m *machineEnvStore) KeyName() string {
	return m.KeyNameValue
}

// Type returns the type value read from env vars.
func (m *machineEnvStore) Type() string {
	return m.MachineTypeValue
}

// FirewallRules returns the firewall rules value read from env vars.
func (m *machineEnvStore) FirewallRules() []string {
	return m.FirewallRulesValue
}

// subnet calculates and returns the subnet id for the current subnetZoneIndex.
func (m *machineEnvStore) subnet() string {
	i := m.subnetZoneIndex % len(m.SubnetsValue)
	return m.SubnetsValue[i]
}

// zone calculates and returns the zone id for the current subnetZoneIndex.
func (m *machineEnvStore) zone() string {
	i := m.subnetZoneIndex % len(m.ZonesValue)
	return m.SubnetsValue[i]
}

// Timeout calculates the time duration in seconds for the current NodeReadyTimeout value.
func (m *machineEnvStore) Timeout() time.Duration {
	return time.Duration(m.NodeReadyTimeout) * time.Second
}

//PollFrequency returns a time duration of 2 seconds.
func (m *machineEnvStore) PollFrequency() time.Duration {
	return 2 * time.Second
}

// SubnetAndZone returns the subnet and zone.
// It performs a round robin operation incrementing the subnetZoneIndex.
func (m *machineEnvStore) SubnetAndZone() (string, string) {
	subnet, zone := m.subnet(), m.zone()
	m.subnetZoneIndex++
	return subnet, zone
}

// Tags creates a set of tags for a certain machine using the given simulation, nodeType and nameSuffix.
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

// NamePrefix returns the name prefix value.
func (m *machineEnvStore) NamePrefix() string {
	return m.NamePrefixValue
}

// Limit returns the maximum amount of machines that can be created.
func (m *machineEnvStore) Limit() int {
	return m.MachinesLimitValue
}

// InitScript returns the script that will be run when the machine gets created.
// TODO: Address this function when implementing the corresponding job.
func (m *machineEnvStore) InitScript() *string {
	return nil
}

// newMachinesStore initializes a new store.Machines implementation using machineEnvStore.
func newMachinesStore() store.Machines {
	var m machineEnvStore
	if err := env.Parse(&m); err != nil {
		panic(err)
	}
	return &m
}
