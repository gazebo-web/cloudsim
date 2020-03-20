package orchestrator

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

type KubernetesMock struct {
	kubernetes.Interface
}

// NewTest creates a new Kubernetes client for testing purposes.
func NewTest() kubernetes.Interface {
	return &KubernetesMock{Interface: fake.NewSimpleClientset()}
}

func (kc *KubernetesMock) SetClientset(cli kubernetes.Interface) {
	kc.Interface = cli
}