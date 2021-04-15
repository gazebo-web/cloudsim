package manager

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	platformFactory "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform/implementations"
)

// Map is the default Manager implementation.
type Map map[Selector]platform.Platform

// Selectors returns a slice with all the available platform selectors.
func (p Map) Selectors() []Selector {
	selectors := make([]Selector, len(p))
	i := 0
	for selector := range p {
		selectors[i] = selector
		i++
	}

	return selectors
}

// Platforms returns a slice with all the available platforms.
func (p Map) Platforms() []platform.Platform {
	platforms := make([]platform.Platform, len(p))
	i := 0
	for _, platform := range p {
		platforms[i] = platform
		i++
	}

	return platforms
}

// Platform returns receives a selector and returns its matching platform or an error if it is not found.
func (p Map) Platform(selector Selector) (platform.Platform, error) {
	platform, ok := p[selector]
	if !ok {
		return nil, ErrPlatformNotFound
	}

	return platform, nil
}

// Platform returns receives a selector and returns its matching platform or an error if it is not found.
func (p Map) set(selector Selector, platform platform.Platform) error {
	// Fail if the platform has already been defined
	_, ok := p[selector]
	if ok {
		return errors.Wrap(ErrPlatformExists, fmt.Sprintf("failed to set platform %s", selector))
	}
	// Register the platform
	p[selector] = platform

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
		if err := m.set(Selector(name), out); err != nil {
			return nil, err
		}
	}

	return m, nil
}
