package manager

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/loader"
	"github.com/gazebo-web/cloudsim/v4/pkg/platform"
	"github.com/gazebo-web/cloudsim/v4/pkg/simulations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/pkg/errors"
	"path/filepath"
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
//   - aws_k8s_us_east_1 - Platform containing AWS and Kubernetes components pointed at us-east-1.
//   - aws_k8s_us_east_2 - Platform containing AWS and Kubernetes components pointed at us-east-2.
//   - gcp_mesos - Platform containing GCE and Apache Mesos components.
//   - azure_swarm - Platform containing Azure and Docker Swarm components.
type Manager interface {
	// Selectors returns a slice with all the available platform selectors.
	Selectors() []string
	// Platforms returns a slice with all the available platforms.
	// The `selector` parameter can be passed to define a target platform. If `selector` is defined and a platform is
	// matched, the matched platform will be the first element in the returned list. If `target` is `nil` or is not
	// found, the elements in the list will be returned in random order.
	Platforms(selector *string) []platform.Platform
	// Platform returns the platform that matches a specific selector, or an error if it is not found.
	Platform(selector string) (platform.Platform, error)
}

// managerConfig defines the platform configuration file structure.
type managerConfig struct {
	// Platforms contains information used to create platforms.
	// A platform will be created for each entry in the map.
	// Keys define platform names.
	// Values must must match the platform factory Config struct.
	Platforms map[string]factory.Config
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
	Logger gz.Logger
}

// loadPlatformConfiguration loads platform configuration files from NewInput.ConfigPath and returns a loaded managerConfig
// value.
//
//	If NewInput.ConfigPath is a directory, all config files within that directory will be loaded.
//	If NewInput.ConfigPath is a file, it will only use that file as the config file.
//	A `name` config value containing the platform name will be added to each platform's factory config fields.
func loadPlatformConfiguration(input *NewInput) (*managerConfig, error) {
	list, err := listConfigFiles(input.ConfigPath)
	if err != nil {
		return nil, err
	}

	mc := managerConfig{
		Platforms: make(map[string]factory.Config),
	}

	list = input.Loader.Filter(list)

	for _, p := range list {
		var config factory.Config

		err = input.Loader.Load(p, &config)
		if err != nil {
			continue
		}

		// Get filename as key for platform map
		file := filepath.Base(p)
		filename := input.Loader.TrimExt(file)

		mc.Platforms[filename] = config
		mc.Platforms[filename].Config["name"] = filename
	}

	return &mc, err
}

// GetSimulationPlatform gets the platform.Platform associated with a simulation.
func GetSimulationPlatform(manager Manager, sim simulations.Simulation) (platform.Platform, error) {
	platformName := sim.GetPlatform()
	if platformName == nil {
		return nil, simulations.ErrSimulationPlatformNotDefined
	}

	return manager.Platform(*platformName)
}
