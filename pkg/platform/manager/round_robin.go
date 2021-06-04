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
	return &roundRobin{
		Cycler:  c,
		Manager: base,
	}, nil
}

// roundRobin is a Manager implementation using round robin technique.
type roundRobin struct {
	cycler.Cycler
	Manager
}

// Platforms returns a slice with all the available manager, but compared to Map.Platforms, it will try to
// return a different platform every time at the index 0 if no selector is passed using round robin.
func (c *roundRobin) Platforms(selector *string) []platform.Platform {
	if selector == nil {
		next := c.Next().(string)
		return c.Platforms(&next)
	}
	return c.Platforms(selector)
}
