package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetConfig returns the rest config for accessing a Kubernetes master.
// It reads the configuration from the default .kube config path.
func GetConfig() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	return config, nil
}

// NewAPI returns a client object to access a kubernetes master.
func NewAPI(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
