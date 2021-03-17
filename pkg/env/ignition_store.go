package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/network"
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

	// VerbosityValue is the IGN_VERBOSE value that will be passed to Pods launched for SubT.
	VerbosityValue string `env:"CLOUDSIM_IGN_VERBOSITY" envDefault:"2"`

	// LogsCopyEnabledValue is the CLOUDSIM_IGN_LOGS_COPY_ENABLED value that will used to define if logs should be copied.
	LogsCopyEnabledValue bool `env:"CLOUDSIM_IGN_LOGS_COPY_ENABLED"`

	// RegionValue is the CLOUDSIM_IGN_REGION value that will determine where to launch simulations.
	RegionValue string `env:"CLOUDSIM_IGN_REGION"`

	// SecretsNameValue is the CLOUDSIM_IGN_SECRETS_NAME value that will used to get credentials for cloud providers.
	SecretsNameValue string `env:"CLOUDSIM_IGN_SECRETS_NAME"`

	// LogsBucketValue is the CLOUDSIM_AWS_GZ_LOGS_BUCKET value that will be used to upload logs.
	LogsBucketValue string `env:"CLOUDSIM_AWS_GZ_LOGS_BUCKET"`

	// WebsocketHostValue is the CLOUDSIM_WEBSOCKET_HOST that will be used as host to connect to simulation's websocket servers.
	WebsocketHostValue string `env:"CLOUDSIM_SUBT_WEBSOCKET_HOST"`

	// DefaultRecipientsValue has the list of emails that should always receive summaries.
	DefaultRecipientsValue []string `env:"CLOUDSIM_IGN_DEFAULT_RECIPIENTS"`

	// DefaultSenderValue is the email address used to send emails.
	DefaultSenderValue string `env:"CLOUDSIM_IGN_DEFAULT_SENDER,required"`
}

// DefaultRecipients returns the list of default summary email recipients.
func (i *ignitionEnvStore) DefaultRecipients() []string {
	return i.DefaultRecipientsValue
}

// DefaultSender returns the default email address used to send emails.
func (i *ignitionEnvStore) DefaultSender() string {
	return i.DefaultSenderValue
}

// LogsBucket returns the bucket to upload simulation logs to.
func (i *ignitionEnvStore) LogsBucket() string {
	return i.LogsBucketValue
}

// GetWebsocketHost returns the host of the websocket address for connecting to simulation websocket servers.
func (i *ignitionEnvStore) GetWebsocketHost() string {
	return i.WebsocketHostValue
}

// GetWebsocketPath returns the path of the websocket address for the given simulation's group id.
func (i *ignitionEnvStore) GetWebsocketPath(groupID simulations.GroupID) string {
	return fmt.Sprintf("simulations/%s", groupID.String())
}

// AccessKeyLabel returns the access key label to get the credentials for a certain cloud provider.
// For AWS, it returns: `aws-access-key-id`
func (i *ignitionEnvStore) AccessKeyLabel() string {
	return "aws-access-key-id"
}

// SecretAccessKeyLabel returns the secret access key label to get the credentials for a certain cloud provider.
// For AWS, it returns: `aws-secret-access-key`
func (i *ignitionEnvStore) SecretAccessKeyLabel() string {
	return "aws-secret-access-key"
}

// LogsCopyEnabled determines if ROS/Gazebo logs should be saved in a bucket or not.
func (i *ignitionEnvStore) LogsCopyEnabled() bool {
	return i.LogsCopyEnabledValue
}

// Region returns the region where to launch a certain simulation.
func (i *ignitionEnvStore) Region() string {
	return i.RegionValue
}

// SecretsName returns the name of the secrets to access credentials for different cloud providers.
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
// If the IP address cannot be obtained, an empty string will be returned.
func (i *ignitionEnvStore) IP() string {
	ip, err := network.GetLocalIPAddressString()
	if err != nil {
		return ""
	}
	return ip
}

// newIgnitionStore initializes a new store.Ignition implementation using ignitionEnvStore.
func newIgnitionStore() (store.Ignition, error) {
	var i ignitionEnvStore
	if err := env.Parse(&i); err != nil {
		return nil, err
	}
	return &i, nil
}
