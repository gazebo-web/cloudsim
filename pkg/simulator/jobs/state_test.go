package jobs

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/actions"
	"github.com/gazebo-web/cloudsim/v4/pkg/platform"
)

// TestState is a type that implements the actions.State interface and related interfaces.
// It is defined here to reduce the amount of duplicated code required to test jobs.
type TestState struct {
	platform platform.Platform
}

// ToStore is a helper method that wraps this state in an Actions store.
func (s *TestState) Platform() platform.Platform {
	return s.platform
}

// ToStore is a helper method that wraps this state in an Actions store.
func (s *TestState) ToStore() actions.Store {
	return actions.NewStore(s)
}
