package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// machineEnvStore is a store.Machines implementation.
// It contains all the information needed by application jobs to launch simulations.
type machineEnvStore struct {
	// InstanceProfileValue is the ARN used to configure EC2 machines.
	InstanceProfileValue string `env:"CLOUDSIM_MACHINES_INSTANCE_PROFILE" envDefault:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker"`

	// KeyNameValue is the name of the SSH key used for a new instance.
	KeyNameValue string `env:"CLOUDSIM_MACHINES_KEY_NAME" envDefault:"ignitionFuel"`

	// MachineTypeValue is the type of instance that will be created.
	MachineTypeValue string `env:"CLOUDSIM_MACHINES_TYPE" envDefault:"g3.4xlarge"`

	// FirewallRulesValue is a set of firewall rules that will be applied to a new instance.
	FirewallRulesValue []string `env:"CLOUDSIM_MACHINES_FIREWALL_RULES" envSeparator:"," envDefault:"sg-0c5c791266694a3ca"`

	// SubnetsValue is a slice of AWS subnet IDs to launch simulations in. (Example: subnet-1270518251)
	SubnetsValue []string `env:"CLOUDSIM_MACHINES_SUBNETS,required" envSeparator:","`

	// ZonesValue is a slice of AWS availability zones to launch simulations in. (Example: us-east-1a)
	ZonesValue []string `env:"CLOUDSIM_MACHINES_ZONES,required" envSeparator:","`

	// MachinesLimitValue is the maximum number of machines that Cloudsim can have running at the same time.
	MachinesLimitValue int `env:"CLOUDSIM_MACHINES_LIMIT" envDefault:"-1"`

	// BaseImageValue is the Amazon Machine Image name that is used as base image for the a new instance.
	BaseImageValue string `env:"CLOUDSIM_MACHINES_BASE_IMAGE" envDefault:"ami-08861f7e7b409ed0c"`

	// NamePrefixValue is the prefix used when naming a new instance.
	NamePrefixValue string `env:"CLOUDSIM_MACHINES_NAME_PREFIX,required" envDefault:"cloudsim-subt-node"`

	// ClusterNameValue contains the name of the cluster EC2 instances will join.
	ClusterNameValue string `env:"CLOUDSIM_MACHINES_CLUSTER_NAME,required"`

	// NodeReadyTimeout is the total amount of time in seconds that the machine creation process will wait.
	NodeReadyTimeout uint `env:"CLOUDSIM_MACHINES_NODE_READY_TIMEOUT_SECONDS" envDefault:"300"`

	// subnetZoneIndex is used as round robin index for setting different subnets and zones to different machines.
	subnetZoneIndex int
}

// BaseImage returns the base image value read from env vars.
// In AWS, the base image is the Amazon Machine Image (AMI).
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

// PollFrequency returns a time duration of 2 seconds.
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
func (m *machineEnvStore) Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []machines.Tag {
	name := fmt.Sprintf("%s-%s-%s", m.NamePrefixValue, simulation.GetGroupID(), nameSuffix)
	clusterKey := fmt.Sprintf("kubernetes.io/cluster/%s", m.ClusterNameValue)
	return []machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"Name":                       name,
				"cloudsim_groupid":           string(simulation.GetGroupID()),
				"CloudsimGroupID":            string(simulation.GetGroupID()),
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
