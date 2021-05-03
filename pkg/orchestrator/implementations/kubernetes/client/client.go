package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// GetConfig returns the rest config for accessing a Kubernetes master.
// It reads the configuration from the default .kube config path.
func GetConfig(path string) (*rest.Config, error) {
	// Use the default config if no config is specified.
	// Check the KUBECONFIG env var
	if path == "" {
		path = os.Getenv("KUBECONFIG")
	}
	// If the path is still empty, default to the user's .kube directory config
	if path == "" {
		path = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	// Get the config
	config, err := clientcmd.BuildConfigFromFlags("", path)
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
