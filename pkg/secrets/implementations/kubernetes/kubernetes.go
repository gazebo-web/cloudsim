package kubernetes

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// kubernetesSecrets uses a kubernetes client to read secrets from the cluster.
type kubernetesSecrets struct {
	client v1.SecretsGetter
}

// Get gets a certain secret with the given name and in the given namespace.
func (s *kubernetesSecrets) Get(ctx context.Context, name, namespace string) (*secrets.Secret, error) {
	sc, err := s.client.Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &secrets.Secret{
		Data: sc.Data,
	}, nil
}

// NewKubernetesSecrets initializes a new Secrets implementation using Kubernetes.
func NewKubernetesSecrets(client v1.SecretsGetter) secrets.Secrets {
	return &kubernetesSecrets{
		client: client,
	}
}
