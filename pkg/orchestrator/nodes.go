package orchestrator

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
)

var (
	// ErrNodesNotReady is returned when the nodes are not ready.
	ErrNodesNotReady = errors.New("nodes are not ready")
)

// Nodes groups a set of methods to register nodes into a cluster.
type Nodes interface {
	WaitForCondition(node Resource, condition Condition) waiter.Waiter
}
