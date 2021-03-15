package secrets

import (
	"context"
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
