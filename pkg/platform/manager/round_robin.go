package manager

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/platform"
	"github.com/gazebo-web/gz-go/v7/cycler"
)

// WithRoundRobin initializes a new Manager using a base platform Manager but iterates over the inner manager using round robin.
func WithRoundRobin(base Manager) (Manager, error) {
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

// Platforms returns a slice with a set of platforms.
// If a selector is passed, the underlying Manager will be in charge of returning the platform slice.
// If no selector is passed, this implementation will try to return a platform based on a round robin algorithm.
func (c *roundRobin) Platforms(selector *string) []platform.Platform {
	if selector == nil {
		next := c.Next().(string)
		return c.Manager.Platforms(&next)
	}
	return c.Manager.Platforms(selector)
}
