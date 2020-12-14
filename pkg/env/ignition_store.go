package env

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// ignitionEnvStore is the implementation of store.Ignition using env vars.
type ignitionEnvStore struct {
	// GazeboServerLogsPathValue is the path inside the container where the `gz-logs` Volume will be mounted.
	// eg. '/tmp/ign'.
	// Important: it is important that gazebo is configured to record its logs to a child folder of the
	// passed mount location (eg. following the above example, '/var/log/gzserver/logs'), otherwise gazebo
	// will 'mkdir' and override the mounted folder.
	// See the "gzserver-container" Container Spec below to see the code.
	GazeboServerLogsPathValue string `env:"CLOUDSIM_IGN_GZSERVER_LOGS_VOLUME_MOUNT_PATH" envDefault:"/tmp/ign"`

	// ROSLogsPathValue is the path inside the ROS container where the ros logs Volume will be mounted.
	ROSLogsPathValue string `env:"CLOUDSIM_IGN_BRIDGE_LOGS_VOLUME_MOUNT_PATH" envDefault:"/home/developer/.ros"`

	// SidecarContainerLogsPathValue is the path inside the sidecar container where the logs volume will be mounted.
	SidecarContainerLogsPathValue string `env:"CLOUDSIM_IGN_SIDECAR_CONTAINER_LOGS_VOLUME_MOUNT_PATH" envDefault:"/tmp/logs"`

	// IgnIPValue is the Cloudsim server's IP address to use when creating NetworkPolicies.
	// See 'docker-entrypoint.sh' script located at the root folder of this project.
	IgnIPValue string `env:"CLOUDSIM_IGN_IP"`

	// VerbosityValue is the IGN_VERBOSE value that will be passed to Pods launched for SubT.
	VerbosityValue string `env:"CLOUDSIM_IGN_VERBOSITY"`

	// LogsCopyEnabledValue is the CLOUDSIM_IGN_LOGS_COPY_ENABLED value that will used to define if logs should be copied.
	LogsCopyEnabledValue bool `env:"CLOUDSIM_IGN_LOGS_COPY_ENABLED"`

	// RegionValue is the CLOUDSIM_IGN_REGION value that will determine where to launch simulations.
	RegionValue string `env:"CLOUDSIM_IGN_REGION"`

	// SecretsNameValue is the CLOUDSIM_IGN_SECRETS_NAME value that will used to get credentials for cloud providers.
	SecretsNameValue string `env:"CLOUDSIM_IGN_SECRETS_NAME"`

	GazeboBucketValue string `env:"CLOUDSIM_IGN_GAZEBO_BUCKET,required"`
}

func (i *ignitionEnvStore) GazeboBucket() string {
	return i.GazeboBucketValue
}

func (i *ignitionEnvStore) AccessKeyLabel() string {
	return "aws-access-key-id"
}

func (i *ignitionEnvStore) SecretAccessKeyLabel() string {
	return "aws-secret-access-key"
}

func (i *ignitionEnvStore) LogsCopyEnabled() bool {
	return i.LogsCopyEnabledValue
}

func (i *ignitionEnvStore) Region() string {
	return i.RegionValue
}

func (i *ignitionEnvStore) SecretsName() string {
	return i.SecretsNameValue
}

// ROSLogsPath returns the path of the logs from bridge containers.
func (i *ignitionEnvStore) ROSLogsPath() string {
	return i.ROSLogsPathValue
}

// SidecarContainerLogsPath returns the path of the logs from sidecar containers.
func (i *ignitionEnvStore) SidecarContainerLogsPath() string {
	return i.SidecarContainerLogsPathValue
}

// GazeboServerLogsPath returns the path of the logs from gazebo server containers.
func (i *ignitionEnvStore) GazeboServerLogsPath() string {
	return i.GazeboServerLogsPathValue
}

// Verbosity returns the level of verbosity that should be used for gazebo.
func (i *ignitionEnvStore) Verbosity() string {
	return i.VerbosityValue
}

// IP returns the Cloudsim server's IP address to use when creating NetworkPolicies.
func (i *ignitionEnvStore) IP() string {
	return i.IgnIPValue
}

// newIgnitionStore initializes a new store.Ignition implementation using ignitionEnvStore.
func newIgnitionStore() store.Ignition {
	var i ignitionEnvStore
	if err := env.Parse(&i); err != nil {
		panic(err)
	}
	return &i
}
