package store

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/creasty/defaults"
	defaulter "gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"time"
)

// machinesStore is a store.Machines implementation.
// It contains all the information needed by application jobs to launch simulations.
type machinesStore struct {
	// InstanceProfileValue is the ARN used to configure EC2 machines.
	InstanceProfileValue string `default:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker" env:"CLOUDSIM_MACHINES_INSTANCE_PROFILE"`

	// KeyNameValue is the name of the SSH key used for a new instance.
	KeyNameValue string `default:"ignitionFuel" env:"CLOUDSIM_MACHINES_KEY_NAME"`

	// MachineTypeValue is the type of instance thGat will be created.
	MachineTypeValue string `default:"g3.4xlarge" env:"CLOUDSIM_MACHINES_TYPE"`

	// FirewallRulesValue is a set of firewall rules that will be applied to a new instance.
	FirewallRulesValue []string `default:"[\"sg-0c5c791266694a3ca\"]" env:"CLOUDSIM_MACHINES_FIREWALL_RULES" envSeparator:","`

	// SubnetsValue is a slice of AWS subnet IDs to launch simulations in. (Example: subnet-1270518251)
	SubnetsValue []string `validate:"required" env:"CLOUDSIM_MACHINES_SUBNETS,required" envSeparator:","`

	// ZonesValue is a slice of AWS availability zones to launch simulations in. (Example: us-east-1a)
	ZonesValue []string `validate:"required" env:"CLOUDSIM_MACHINES_ZONES,required" envSeparator:","`

	// MachinesLimitValue is the maximum number of machines that Cloudsim can have running at the same time.
	MachinesLimitValue int `default:"-1" env:"CLOUDSIM_MACHINES_LIMIT"`

	// BaseImageValue is the Amazon Machine Image name that is used as base image for the a new instance.
	BaseImageValue string `default:"ami-08861f7e7b409ed0c" env:"CLOUDSIM_MACHINES_BASE_IMAGE"`

	// NamePrefixValue is the prefix used when naming a new instance.
	NamePrefixValue string `validate:"required" default:"cloudsim-subt-node" env:"CLOUDSIM_MACHINES_NAME_PREFIX,required"`

	// ClusterNameValue contains the name of the cluster EC2 instances will join.
	ClusterNameValue string `validate:"required" env:"CLOUDSIM_MACHINES_CLUSTER_NAME,required"`

	// NodeReadyTimeout is the total amount of time in seconds that the machine creation process will wait.
	NodeReadyTimeout uint `default:"300" env:"CLOUDSIM_MACHINES_NODE_READY_TIMEOUT_SECONDS"`

	// subnetZoneIndex is used as round robin index for setting different subnets and zones to different machines.
	subnetZoneIndex int
}

// SetDefaults sets default values for the store.
func (m *machinesStore) SetDefaults() error {
	defaults.MustSet(m)
	return nil
}

// ClusterName returns the cluster name.
// In AWS: It returns the EKS cluster name.
func (m *machinesStore) ClusterName() string {
	return m.ClusterNameValue
}

// BaseImage returns the base image value read from env vars.
// In AWS, the base image is the Amazon Machine Image (AMI).
func (m *machinesStore) BaseImage() string {
	return m.BaseImageValue
}

// InstanceProfile returns the instance profile value read from env vars.
func (m *machinesStore) InstanceProfile() *string {
	return &m.InstanceProfileValue
}

// KeyName returns the key name value read from env vars.
func (m *machinesStore) KeyName() string {
	return m.KeyNameValue
}

// Type returns the type value read from env vars.
func (m *machinesStore) Type() string {
	return m.MachineTypeValue
}

// FirewallRules returns the firewall rules value read from env vars.
func (m *machinesStore) FirewallRules() []string {
	return m.FirewallRulesValue
}

// subnet calculates and returns the subnet id for the current subnetZoneIndex.
func (m *machinesStore) subnet() string {
	i := m.subnetZoneIndex % len(m.SubnetsValue)
	return m.SubnetsValue[i]
}

// zone calculates and returns the zone id for the current subnetZoneIndex.
func (m *machinesStore) zone() string {
	i := m.subnetZoneIndex % len(m.ZonesValue)
	return m.ZonesValue[i]
}

// Timeout calculates the time duration in seconds for the current NodeReadyTimeout value.
func (m *machinesStore) Timeout() time.Duration {
	return time.Duration(m.NodeReadyTimeout) * time.Second
}

// PollFrequency returns a time duration of 2 seconds.
func (m *machinesStore) PollFrequency() time.Duration {
	return 2 * time.Second
}

// SubnetAndZone returns the subnet and zone.
// It performs a round robin operation incrementing the subnetZoneIndex.
func (m *machinesStore) SubnetAndZone() (string, string) {
	subnet, zone := m.subnet(), m.zone()
	m.subnetZoneIndex++
	return subnet, zone
}

// Tags creates a set of tags for a certain machine using the given simulation, nodeType and nameSuffix.
func (m *machinesStore) Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []machines.Tag {
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
func (m *machinesStore) NamePrefix() string {
	return m.NamePrefixValue
}

// Limit returns the maximum amount of machines that can be created.
func (m *machinesStore) Limit() int {
	return m.MachinesLimitValue
}

// newMachinesStoreFromEnvVars initializes a new store.Machines implementation using environment variables to
// configure a machinesStore object.
func newMachinesStoreFromEnvVars() (storepkg.Machines, error) {
	// Load store from env vars
	var m machinesStore
	if err := env.Parse(&m); err != nil {
		return nil, err
	}
	// Set default values
	if err := defaulter.SetDefaults(&m); err != nil {
		return nil, err
	}
	// Validate values
	if err := validate.Validate(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
