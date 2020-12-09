package secrets

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Secret represents a set of secret data.
type Secret struct {
	// Data contains the actual secret information.
	Data map[string][]byte
}

// Secrets has a set of methods to get secrets for a certain name and namespace.
type Secrets interface {
	Get(ctx context.Context, name, namespace string) (*Secret, error)
}

// kubernetesSecrets uses a kubernetes client to read secrets from the cluster.
type kubernetesSecrets struct {
	client v1.SecretsGetter
}

// Get gets a certain secret with the given name and in the given namespace.
func (s *kubernetesSecrets) Get(ctx context.Context, name, namespace string) (*Secret, error) {
	sc, err := s.client.Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &Secret{
		Data: sc.Data,
	}, nil
}

// NewKubernetesSecrets initializes a new Secrets implementation using Kubernetes.
func NewKubernetesSecrets(client v1.SecretsGetter) Secrets {
	return &kubernetesSecrets{
		client: client,
	}
}
