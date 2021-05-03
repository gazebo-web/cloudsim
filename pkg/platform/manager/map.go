package manager

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	platformFactory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform/implementations"
)

// Map is the default Manager implementation.
type Map map[string]platform.Platform

// Selectors returns a slice with all the available platform selectors.
func (m Map) Selectors() []string {
	selectors := make([]string, len(m))
	i := 0
	for selector := range m {
		selectors[i] = selector
		i++
	}

	return selectors
}

// Platforms returns a slice with all the available platforms.
func (m Map) Platforms(selector *string) []platform.Platform {
	var platforms []platform.Platform
	for key, p := range m {
		// If a selector was provided and is matched, place at the front of the slice
		if selector != nil && *selector == key {
			platforms = append([]platform.Platform{p}, platforms...)
		} else {
			platforms = append(platforms, p)
		}
	}

	return platforms
}

// Platform receives a selector and returns its matching platform or an error if it is not found.
func (m Map) Platform(selector string) (platform.Platform, error) {
	platform, ok := m[selector]
	if !ok {
		return nil, ErrPlatformNotFound
	}

	return platform, nil
}

// Platform returns receives a selector and returns its matching platform or an error if it is not found.
func (m Map) set(selector string, platform platform.Platform) error {
	// Fail if the platform has already been defined
	_, ok := m[selector]
	if ok {
		return errors.Wrap(ErrPlatformExists, fmt.Sprintf("failed to set platform %s", selector))
	}
	// Register the platform
	m[selector] = platform

	return nil
}

// NewMapFromConfig loads a platformMap of platforms from a configuration file and returns a platform Map containing
// the platforms.
func NewMapFromConfig(input *NewInput) (Manager, error) {
	if input == nil {
		return nil, ErrInvalidNewInput
	}

	// Load config
	fileConfig, err := loadPlatformConfiguration(input)
	if err != nil {
		return nil, err
	}

	// Prepare dependencies
	dependencies := factory.Dependencies{
		"logger": input.Logger,
	}

	// Create and load map
	m := make(Map, 0)
	for name, config := range fileConfig.Platforms {
		// Create platform
		var out platform.Platform
		if err := platformFactory.Factory.New(config, dependencies, &out); err != nil {
			return nil, err
		}

		// Add platform to map
		if err := m.set(name, out); err != nil {
			return nil, err
		}
	}

	return m, nil
}
