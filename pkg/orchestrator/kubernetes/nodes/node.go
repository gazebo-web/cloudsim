package nodes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

// Node is a Kubernetes node. It extends the generic orchestrator.Resource interface.
type Node interface {
	orchestrator.Resource
}

// node is a Node implementation that contains the basic information to identify a node in a Kubernetes cluster.
type node struct {
	name      string
	selector  string
	namespace string
}

// Name returns the name of the node.
func (n *node) Name() string {
	return n.name
}

// Selector returns the selector of the node.
func (n *node) Selector() string {
	return n.selector
}

// Namespace returns the namespace of the node.
func (n *node) Namespace() string {
	return n.namespace
}

// NewNode returns a new Node implementation using node.
func NewNode(name string, namespace string, selector string) Node {
	return &node{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
