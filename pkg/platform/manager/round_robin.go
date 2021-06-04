package manager

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cycler"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// NewRoundRobin initializes a new manager using a base platform map but iterates over the map using round robin.
func NewRoundRobin(platforms Map) (Manager, error) {
	c, err := cycler.NewCycler(platforms.Selectors())
	if err != nil {
		return nil, err
	}
	return &RoundRobin{
		cycler:    c,
		platforms: platforms,
	}, nil
}

// RoundRobin is a Manager implementation using round robin technique.
type RoundRobin struct {
	cycler    cycler.Cycler
	platforms Map
}

// Selectors returns a slice with all the available platform selectors.
func (c *RoundRobin) Selectors() []string {
	return c.platforms.Selectors()
}

// Platforms returns a slice with all the available platforms, but compared to Map.Platforms, it will try to
// return a different platform every time if no selector is passed using round robin.
func (c *RoundRobin) Platforms(selector *string) []platform.Platform {
	if selector == nil {
		next := c.cycler.Next().(string)
		return c.platforms.Platforms(&next)
	}
	return c.platforms.Platforms(selector)
}

// Platform receives a selector and returns its matching platform or an error if it is not found.
func (c *RoundRobin) Platform(selector string) (platform.Platform, error) {
	return c.platforms.Platform(selector)
}
