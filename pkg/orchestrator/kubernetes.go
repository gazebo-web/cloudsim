package orchestrator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Kubernetes struct {
	kubernetes.Interface
}

// New creates a new Kubernetes client to access a kubernetes master.
func New() *Kubernetes {
	kcli, err := NewClientset()
	if err != nil {
		log.Fatalf("[KUBERNETES] Error trying to create a client to kubernetes %+v\n", err)
	}
	k := Kubernetes{kcli}
	return &k
}

// NewClientset creates a new Kubernetes Clientset from the kubeconfig file.
// Note that this kube client assumes there is a kubernetes configuration in the
// server's ~/.kube/config file. That config is used to connect to the kubernetes
// master.
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

func (kc *Kubernetes) SetClientset(cli kubernetes.Interface) {
	kc.Interface = cli
}

// MakeListOptions returns a ListOptions object for an array of labels.
func MakeListOptions(labels ...string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	}
}