package manager

import (
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/loader"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"os"
	"path"
)

var (
	// ErrPlatformExists is returned when a platform for a provided selector exists. Used to prevent replacing
	// platforms.
	ErrPlatformExists = errors.New("platform already exists")
	// ErrPlatformNotFound is returned when a platform for a provided selector cannot be found.
	ErrPlatformNotFound = errors.New("platform not found")
	// ErrInvalidNewInput is returned when an invalid input is provided to a Manager implementation creation
	// function.
	ErrInvalidNewInput = errors.New("invalid NewMap input")
)

// Selector is used to uniquely identify a Platform.
type Selector string

// Manager manages a platformMap of platforms that can be used by Cloudsim.
// Implementations must be simple platform containers and must not be aware of the feature sets, differences (if any) or
// implementation details of the available platforms they contain. This is by design, in order to give applications the
// flexibility to manage sets of platforms as they see fit.
// All contained platforms are completely independent from each other. As an example, you can have a platformMap of platforms
// that use the same internal components (e.g. AWS+Kubernetes), but that are configured to point to different
// service regions. At the same time, it is possible to have additional platforms making use of entirely different
// components (e.g. GCP+Mesos, Azure+Swarm) co-exist with the platformMap of AWS+Kubernetes platforms.
// Every platform is uniquely identified by a Selector, a user-defined identifier. To make use of a specific platform,
// the target platform's Selector is passed to the Manager implementation through the Platform method.
// An example of a selector platformMap:
//   * aws_k8s_us_east_1 - Platform containing AWS and Kubernetes components pointed at us-east-1.
//   * aws_k8s_us_east_2 - Platform containing AWS and Kubernetes components pointed at us-east-2.
//   * gcp_mesos - Platform containing GCE and Apache Mesos components.
//   * azure_swarm - Platform containing Azure and Docker Swarm components.
type Manager interface {
	// Selectors returns a slice with all the available platform selectors.
	Selectors() []Selector
	// Platforms returns a slice with all the available platforms.
	Platforms() []platform.Platform
	// Platform returns the platform that matches a specific selector, or an error if it is not found.
	GetPlatform(selector Selector) (platform.Platform, error)
}

// managerConfig defines the platform configuration file structure.
type managerConfig struct {
	// Platforms contains information used to create platforms.
	// A platform will be created for each entry in the map.
	// Keys define platform names.
	// Values must must match the platform factory Config struct.
	Platforms map[string]*factory.Config
}

// NewInput contains common information necessary to create a new manager implementation instance.
// Manager implementations should use or embed this structure to request input data.
type NewInput struct {
	// ConfigPath contains the path to the platforms configuration file.
	// If this field is an empty string, it will default to the `config.yaml` file in the directory Cloudsim is running
	// in.
	ConfigPath string
	// Loader used to load the configuration file
	Loader loader.Loader
	// Logger used to configure platforms.
	Logger ign.Logger
}

// loadPlatformConfiguration loads a platform configuration file and returns a loaded managerConfig value.
// A `name` config value containing the platform name will be added to each platform's factory config fields.
func loadPlatformConfiguration(input *NewInput) (*managerConfig, error) {
	// Set default path if not defined
	if input.ConfigPath == "" {
		var err error
		var cwd string
		if cwd, err = os.Getwd(); err != nil {
			return nil, err
		}
		input.ConfigPath = path.Join(cwd, "config.yaml")
	}

	// Load platform configurations
	config := &managerConfig{}
	err := input.Loader.Load(input.ConfigPath, config)

	// Append each platform name to its set of values
	for name, config := range config.Platforms {
		config.Config["name"] = name
	}

	return config, err
}

// GetSimulationPlatform gets the platform.Platform associated with a simulation.
func GetSimulationPlatform(manager Manager, sim simulations.Simulation) (platform.Platform, error) {
	platformName := sim.GetPlatform()
	if platformName == nil {
		return nil, simulations.ErrSimulationPlatformNotDefined
	}

	return manager.GetPlatform(Selector(*platformName))
}
