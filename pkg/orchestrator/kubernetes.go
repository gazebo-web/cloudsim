package orchestrator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

// NewClient creates a new Kubernetes client.
func NewClient() kubernetes.Interface {
	kcli, err := NewClientset()
	if err != nil {
		log.Fatalf("[KUBERNETES] Error trying to create a client to kubernetes %+v\n", err)
	}
	return kcli
}

// NewTestClient creates a new Kubernetes client for testing purposes.
func NewTestClient() kubernetes.Interface {
	return &MockClientset{Interface: fake.NewSimpleClientset()}
}

// NewClientset creates a new Kubernetes Clientset from the kubeconfig file.
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := NewConfig(nil)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// NewConfig returns the kubernetes config file in the specified path.
// If no path is provided (i.e. nil), then the configuration in ~/.kube/config
// is returned.
func NewConfig(kubeconfig *string) (*rest.Config, error) {
	if kubeconfig == nil {
		kubeconfig = tools.Sptr(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}