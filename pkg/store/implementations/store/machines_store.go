package store

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
	"time"
)

// machinesStore is a store.Machines implementation.
// It contains all the information needed by application jobs to launch simulations.
type machinesStore struct {
	// InstanceProfileValue is the ARN used to configure EC2 machines.
	InstanceProfileValue string `validate:"required" default:"arn:aws:iam::200670743174:instance-profile/aws-eks-role-cloudsim-worker" env:"CLOUDSIM_MACHINES_INSTANCE_PROFILE"`

	// KeyNameValue is the name of the SSH key used for a new instance.
	// This key must be registered in the machines provider.
	KeyNameValue string `validate:"required" default:"ignitionFuel" env:"CLOUDSIM_MACHINES_KEY_NAME"`

	// MachineTypeValue is the type of instance that will be used to request simulation instances.
	MachineTypeValue string `default:"g3.4xlarge" env:"CLOUDSIM_MACHINES_TYPE"`

	// MachineSidecarTypeValue is the type of instance that will be used to request sidecar instances .
	MachineSidecarTypeValue string `default:"c5.4xlarge" env:"CLOUDSIM_MACHINES_SIDECAR_TYPE"`

	// FirewallRulesValue is a set of firewall rules that will be applied to a new instance.
	FirewallRulesValue []string `default:"[\"sg-0c5c791266694a3ca\"]" env:"CLOUDSIM_MACHINES_FIREWALL_RULES" envSeparator:","`

	// MachinesLimitValue is the maximum number of machines that Cloudsim can have running at the same time.
	MachinesLimitValue int `default:"-1" env:"CLOUDSIM_MACHINES_LIMIT"`

	// BaseImageValue is the Amazon Machine Image name that is used as base image for a new instance.
	// This is usually used by machines that need to be launched alongside simulation nodes.
	BaseImageValue string `default:"ami-08861f7e7b409ed0c" env:"CLOUDSIM_MACHINES_BASE_IMAGE"`

	// BaseImageGPUValue is the Amazon Machine Image name that is used as base image for a new simulation instance.
	// This image has support for GPU and X Server.
	BaseImageGPUValue string `default:"ami-08861f7e7b409ed0c" env:"CLOUDSIM_MACHINES_BASE_GPU_IMAGE"`

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
	return defaults.SetStructValues(m)
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

// BaseImageGPU returns the base gpu image value read from env vars.
// In AWS, the base image is the Amazon Machine Image (AMI).
func (m *machinesStore) BaseImageGPU() string {
	return m.BaseImageGPUValue
}

// InstanceProfile returns the instance profile value read from env vars.
func (m *machinesStore) InstanceProfile() *string {
	return &m.InstanceProfileValue
}

// KeyName returns the key name value read from env vars.
func (m *machinesStore) KeyName() string {
	return m.KeyNameValue
}

// Type returns the machine type value read from env vars.
func (m *machinesStore) Type() string {
	return m.MachineTypeValue
}

// SidecarType returns the sidecar machine type value read from env vars.
func (m *machinesStore) SidecarType() string {
	return m.MachineSidecarTypeValue
}

// FirewallRules returns the firewall rules value read from env vars.
func (m *machinesStore) FirewallRules() []string {
	return m.FirewallRulesValue
}

// Timeout calculates the time duration in seconds for the current NodeReadyTimeout value.
func (m *machinesStore) Timeout() time.Duration {
	return time.Duration(m.NodeReadyTimeout) * time.Second
}

// PollFrequency returns a time duration of 2 seconds.
func (m *machinesStore) PollFrequency() time.Duration {
	return 2 * time.Second
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
	if err := defaults.SetValues(&m); err != nil {
		return nil, err
	}
	// Validate values
	if err := validate.Validate(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
