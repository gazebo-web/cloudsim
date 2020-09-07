package env

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
)

// ignitionEnvStore is the implementation of store.Ignition using env vars.
type ignitionEnvStore struct {
	GazeboServerLogsPathValue     string `env:"CLOUDSIM_IGN_GZSERVER_LOGS_VOLUME_MOUNT_PATH" envDefault:"/tmp/ign"`
	IgnIPValue                    string `env:"CLOUDSIM_IGN_IP"`
	VerbosityValue                string `env:"CLOUDSIM_IGN_VERBOSITY"`
	ROSLogsPathValue              string `env:"CLOUDSIM_IGN_BRIDGE_LOGS_VOLUME_MOUNT_PATH" envDefault:"/home/developer/.ros"`
	SidecarContainerLogsPathValue string `env:"CLOUDSIM_IGN_SIDECAR_CONTAINER_LOGS_VOLUME_MOUNT_PATH" envDefault:"/tmp/logs"`
}

// ROSLogsPath returns the path of the logs from bridge containers.
func (i *ignitionEnvStore) ROSLogsPath() string {
	return i.ROSLogsPathValue
}

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
