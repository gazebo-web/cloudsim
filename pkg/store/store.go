package store

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// Store provides a set of components to store data that needs to be accessed by different services.
type Store interface {
	// Machines provides access to a set of configurations for creating machines.
	Machines() Machines
	// Orchestrator provides access to a set of configurations for cluster management.
	Orchestrator() Orchestrator
	// Ignition provides access to a set of common cloudsim configurations.
	Ignition() Ignition
}

// Machines provides different information for creating machines.
type Machines interface {
	// InstanceProfile returns the machine instance profile that should be used when creating a machine.
	// It returns nil if the default instance profile should be used.
	InstanceProfile() *string

	// KeyName returns the name of the SSH private key that should be used when creating a machine.
	KeyName() string

	// Type returns the machine type that should be used when creating a machine.
	Type() string

	// FirewallRules returns the list of rules that should be applied to the created machine.
	FirewallRules() []string

	// SubnetAndZone returns the subnet and the zone that the created machine should be configured in.
	SubnetAndZone() (string, string)

	// Tags returns a set of tags that will be set to the machine.
	Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []cloud.Tag

	// InitScript returns the script that will be run when creating the machine.
	InitScript() *string

	// BaseImage returns the base image that will be used when creating the machine.
	BaseImage() string

	// Timeout returns the maximum amount of time that a job should wait until a machine is created.
	// Timeout is usually used with PollFrequency.
	Timeout() time.Duration

	// PollFrequency returns the interval of time that a job should wait until performing another request to create machines.
	// PollFrequency is usually used with Timeout.
	PollFrequency() time.Duration

	// Limit returns the maximum limit of machines that can be created.
	Limit() int

	// NamePrefix returns the name prefix that should be used when creating a machine.
	NamePrefix() string
}

// Orchestrator provides different information to manage a cluster.
type Orchestrator interface {
	// Namespace returns the base namespace that should be used for simulations.
	Namespace() string

	// TerminationGracePeriod duration that pods need to terminate gracefully.
	TerminationGracePeriod() time.Duration

	// Nameservers returns a slice of the nameservers used to expose simulations to the internet.
	Nameservers() []string

	// IngressNamespace returns the namespace where the ingress is running.
	IngressNamespace() string

	// IngressName returns the ingress name used by Gloo.
	IngressName() string

	// IngressHost returns a FQDN used to route traffic to cloudsim instances.
	IngressHost() string
}

// Ignition provides general information about cloudsim and ignition gazebo.
type Ignition interface {
	// IP returns the current server's ip.
	IP() string

	// GazeboServerLogsPath returns the path of the logs from gazebo server containers.
	GazeboServerLogsPath() string

	// ROSLogsPath returns the path of the logs from bridge containers.
	ROSLogsPath() string

	// SidecarContainerLogsPath returns the path of the logs from sidecar containers.
	SidecarContainerLogsPath() string

	// Verbosity returns the level of verbosity that should be used for gazebo.
	Verbosity() string

	// LogsCopyEnabled determines if ROS logs should be saved in the storage.
	LogsCopyEnabled() bool

	// Region returns the region where to launch a certain simulation.
	Region() string

	// SecretsName returns the name of the secrets to access credentials for different cloud providers.
	SecretsName() string

	// AccessKeyLabel returns the access key label to get the credentials for a certain cloud provider.
	// For AWS, it returns: `aws-access-key-id`
	AccessKeyLabel() string

	// SecretAccessKeyLabel returns the secret access key label to get the credentials for a certain cloud provider.
	// For AWS, it returns: `aws-secret-access-key`
	SecretAccessKeyLabel() string

	GetWebsocketHost() string

	GetWebsocketPath(groupID simulations.GroupID) string
}
