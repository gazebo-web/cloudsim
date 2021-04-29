package simulations

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// GetKubernetesConfig returns the kubernetes config file in the specified path.
// If no path is provided (i.e. nil), then the configuration in ~/.kube/config
// is returned.
func GetKubernetesConfig(kubeconfig *string) (*restclient.Config, error) {
	if kubeconfig == nil {
		kubeconfig = sptr(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetKubernetesClient returns a client object to access a kubernetes master.
// Note that this kube client assumes there is a kubernetes configuration in the
// server's ~/.kube/config file. That config is used to connect to the kubernetes
// master.
func GetKubernetesClient(kubeconfig *string) (*kubernetes.Clientset, error) {
	config, err := GetKubernetesConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
