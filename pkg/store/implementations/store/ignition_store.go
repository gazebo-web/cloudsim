package store

import (
	"fmt"
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// ignitionStore is the implementation of store.Ignition using env vars.
type ignitionStore struct {
	// GazeboServerLogsPathValue is the path inside the container where the `gz-logs` Volume will be mounted.
	// eg. '/tmp/ign'.
	// Important: it is important that gazebo is configured to record its logs to a child folder of the
	// passed mount location (eg. following the above example, '/var/log/gzserver/logs'), otherwise gazebo
	// will 'mkdir' and override the mounted folder.
	// See the "gzserver-container" Container Spec below to see the code.
	GazeboServerLogsPathValue string `default:"/tmp/ign" env:"CLOUDSIM_IGN_GZSERVER_LOGS_VOLUME_MOUNT_PATH"`

	// ROSLogsPathValue is the path inside the ROS container where the ros logs Volume will be mounted.
	ROSLogsPathValue string `default:"/home/developer/.ros" env:"CLOUDSIM_IGN_BRIDGE_LOGS_VOLUME_MOUNT_PATH"`

	// SidecarContainerLogsPathValue is the path inside the sidecar container where the logs volume will be mounted.
	SidecarContainerLogsPathValue string `default:"/tmp/logs" env:"CLOUDSIM_IGN_SIDECAR_CONTAINER_LOGS_VOLUME_MOUNT_PATH"`

	// IgnIPValue is the Cloudsim server's IP address to use when creating NetworkPolicies.
	// See 'docker-entrypoint.sh' script located at the root folder of this project.
	IgnIPValue string `env:"CLOUDSIM_IGN_IP"`

	// VerbosityValue is the IGN_VERBOSE value that will be passed to Pods launched for SubT.
	VerbosityValue string `default:"2" env:"CLOUDSIM_IGN_VERBOSITY"`

	// LogsCopyEnabledValue is the CLOUDSIM_IGN_LOGS_COPY_ENABLED value that will used to define if logs should be copied.
	LogsCopyEnabledValue bool `env:"CLOUDSIM_IGN_LOGS_COPY_ENABLED"`

	// RegionValue is the CLOUDSIM_IGN_REGION value that will determine where to launch simulations.
	RegionValue string `env:"CLOUDSIM_IGN_REGION"`

	// SecretsNameValue is the CLOUDSIM_IGN_SECRETS_NAME value that will used to get credentials for cloud providers.
	SecretsNameValue string `env:"CLOUDSIM_IGN_SECRETS_NAME"`

	// LogsBucketValue is the CLOUDSIM_AWS_GZ_LOGS_BUCKET value that will be used to upload logs.
	LogsBucketValue string `env:"CLOUDSIM_AWS_GZ_LOGS_BUCKET"`

	// DefaultRecipientsValue has the list of emails that should always receive summaries.
	DefaultRecipientsValue []string `env:"CLOUDSIM_IGN_DEFAULT_RECIPIENTS"`

	// DefaultSenderValue is the email address used to send emails.
	DefaultSenderValue string `validate:"required" env:"CLOUDSIM_IGN_DEFAULT_SENDER"`

	// WebsocketHostValue is the CLOUDSIM_WEBSOCKET_HOST that will be used as host to connect to simulation's websocket servers.
	WebsocketHostValue string `env:"CLOUDSIM_SUBT_WEBSOCKET_HOST"`
}

// LogsBucket returns the bucket to upload simulation logs to.
func (i *ignitionStore) LogsBucket() string {
	return i.LogsBucketValue
}

// DefaultRecipients returns the list of default summary email recipients.
func (i *ignitionStore) DefaultRecipients() []string {
	return i.DefaultRecipientsValue
}

// DefaultSender returns the default email address used to send emails.
func (i *ignitionStore) DefaultSender() string {
	return i.DefaultSenderValue
}

// GetWebsocketHost returns the host of the websocket address for connecting to simulation websocket servers.
func (i *ignitionStore) GetWebsocketHost() string {
	return i.WebsocketHostValue
}

// GetWebsocketPath returns the path of the websocket address for the given simulation's group id.
func (i *ignitionStore) GetWebsocketPath(groupID simulations.GroupID) string {
	return fmt.Sprintf("simulations/%s", groupID.String())
}

// AccessKeyLabel returns the access key label to get the credentials for a certain cloud provider.
// For AWS, it returns: `aws-access-key-id`
func (i *ignitionStore) AccessKeyLabel() string {
	return "aws-access-key-id"
}

// SecretAccessKeyLabel returns the secret access key label to get the credentials for a certain cloud provider.
// For AWS, it returns: `aws-secret-access-key`
func (i *ignitionStore) SecretAccessKeyLabel() string {
	return "aws-secret-access-key"
}

// LogsCopyEnabled determines if ROS/Gazebo logs should be saved in a bucket or not.
func (i *ignitionStore) LogsCopyEnabled() bool {
	return i.LogsCopyEnabledValue
}

// Region returns the region where to launch a certain simulation.
func (i *ignitionStore) Region() string {
	return i.RegionValue
}

// SecretsName returns the name of the secrets to access credentials for different cloud providers.
func (i *ignitionStore) SecretsName() string {
	return i.SecretsNameValue
}

// ROSLogsPath returns the path of the logs from bridge containers.
func (i *ignitionStore) ROSLogsPath() string {
	return i.ROSLogsPathValue
}

// SidecarContainerLogsPath returns the path of the logs from sidecar containers.
func (i *ignitionStore) SidecarContainerLogsPath() string {
	return i.SidecarContainerLogsPathValue
}

// GazeboServerLogsPath returns the path of the logs from gazebo server containers.
func (i *ignitionStore) GazeboServerLogsPath() string {
	return i.GazeboServerLogsPathValue
}

// Verbosity returns the level of verbosity that should be used for gazebo.
func (i *ignitionStore) Verbosity() string {
	return i.VerbosityValue
}

// IP returns the Cloudsim server's IP address to use when creating NetworkPolicies.
func (i *ignitionStore) IP() string {
	return i.IgnIPValue
}

// newIgnitionStoreFromEnvVars initializes a new store.Ignition implementation using environment variables to
// configure an ignitionStore object.
func newIgnitionStoreFromEnvVars() (storepkg.Ignition, error) {
	var i ignitionStore
	if err := env.Parse(&i); err != nil {
		return nil, err
	}
	return &i, nil
}
