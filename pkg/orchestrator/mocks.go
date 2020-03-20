package orchestrator

import "k8s.io/client-go/kubernetes"

type MockClientset struct {
	kubernetes.Interface
}