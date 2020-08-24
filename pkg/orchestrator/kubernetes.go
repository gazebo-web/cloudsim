package orchestrator

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Kubernetes interface {
	kubernetes.Interface
	SetClientset(cli kubernetes.Interface)
	Namespace() string
	KubernetesNodes
	KubernetesPods
}

type KubernetesNodes interface {
	NodeWaitForReady(ctx context.Context, namespace string, groupIDLabel string, timeout time.Duration) error
	NodeWaitToMatchCondition(ctx context.Context, namespace string, opts metav1.ListOptions, timeout time.Duration) error
}

type KubernetesPods interface {
	PodExec(ctx context.Context, namespace string, podName string, container string, command []string, options *remotecommand.StreamOptions) (opts *remotecommand.StreamOptions, err error)
	PodGetLog(ctx context.Context, namespace string, podName string, container string, lines int64) (log *string, err error)
	GetAllPods(label *string) (Pods, error)
	PodCreateExecErrorMessage(errorMsg string, options *remotecommand.StreamOptions) string
	PodWaitForReadyCondition(ctx context.Context, c kubernetes.Interface, namespace string, groupIDLabel string, timeout time.Duration) error
	PodWaitToMatchCondition(ctx context.Context, namespace string, opts metav1.ListOptions, condStr string, timeout time.Duration, condition PodCondition) error
}

// k8s wraps the k8s CLI.
type k8s struct {
	kubernetes.Interface
}

// New creates a new k8s client to access a kubernetes master.
func New() Kubernetes {
	kcli, err := NewClientset()
	if err != nil {
		log.Fatalf("[KUBERNETES] Error trying to create a client to kubernetes %+v\n", err)
	}
	k := k8s{kcli}
	return &k
}

// NewClientset creates a new k8s Clientset from the kubeconfig file.
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

// SetClientset assigns the given cli to the internal Interface
func (kc *k8s) SetClientset(cli kubernetes.Interface) {
	kc.Interface = cli
}

// Namespace returns the default k8s namespace.
func (kc *k8s) Namespace() string {
	return v1.NamespaceDefault
}

// MakeListOptions returns a ListOptions object for an array of labels.
func MakeListOptions(labels ...string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	}
}
