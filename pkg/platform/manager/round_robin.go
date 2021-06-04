package manager

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cycler"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// WithRoundRobin initializes a new Manager using a base platform Manager but iterates over the inner manager using round robin.
func WithRoundRobin(base Manager, err error) (Manager, error) {
	if err != nil {
		return nil, err
	}
	c, err := cycler.NewCyclerFromSlice(base.Selectors())
	if err != nil {
		return nil, err
	}
	return &RoundRobin{
		iterator: c,
		manager:  base,
	}, nil
}

// RoundRobin is a Manager implementation using round robin technique.
type RoundRobin struct {
	iterator cycler.Cycler
	manager  Manager
}

// Selectors returns a slice with all the available platform selectors.
func (c *RoundRobin) Selectors() []string {
	return c.manager.Selectors()
}

// Platforms returns a slice with all the available manager, but compared to Map.Platforms, it will try to
// return a different platform every time at the index 0 if no selector is passed using round robin.
func (c *RoundRobin) Platforms(selector *string) []platform.Platform {
	if selector == nil {
		next := c.iterator.Next().(string)
		return c.manager.Platforms(&next)
	}
	return c.manager.Platforms(selector)
}

// Platform receives a selector and returns its matching platform or an error if it is not found.
func (c *RoundRobin) Platform(selector string) (platform.Platform, error) {
	return c.manager.Platform(selector)
}
