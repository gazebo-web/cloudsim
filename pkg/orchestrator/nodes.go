package orchestrator

import "errors"

var (
	// ErrNodesNotReady is returned when the nodes are not ready.
	ErrNodesNotReady = errors.New("nodes are not ready")
)

// NodeManager groups a set of methods to register nodes into a cluster.
type NodeManager interface {
	Condition(node Resource, condition Condition) Waiter
}
