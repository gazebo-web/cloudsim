package orchestrator

// NodeManager groups a set of methods to register nodes into a cluster.
type NodeManager interface {
	Condition(node Resource, condition Condition) Waiter
}
