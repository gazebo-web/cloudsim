package orchestrator

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// KubernetesMock wraps the Kubernetes CLi for testing purposes.
type KubernetesMock struct {
	kubernetes.Interface
}

// NewTest creates a new Kubernetes client for testing purposes.
func NewTest() kubernetes.Interface {
	return &KubernetesMock{Interface: fake.NewSimpleClientset()}
}

// SetClientset assigns the given cli to the internal Interface
func (kc *KubernetesMock) SetClientset(cli kubernetes.Interface) {
	kc.Interface = cli
}
