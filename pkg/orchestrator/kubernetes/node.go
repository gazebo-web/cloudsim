package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
)

type Node interface {
	orchestrator.Resource
}

type node struct {
	name      string
	selector  string
	namespace string
}

func (n node) Name() string {
	return n.name
}

func (n node) Selector() string {
	return n.selector
}

func (n node) Namespace() string {
	return n.namespace
}

func NewNode(name string, namespace string, selector string) Node {
	return &node{
		name:      name,
		namespace: namespace,
		selector:  selector,
	}
}
