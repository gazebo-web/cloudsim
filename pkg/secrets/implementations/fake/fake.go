package fake

import (
	"context"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
)

// Fake is a fake Secrets implementation.
type Fake struct {
	*mock.Mock
}

// Get mocks the Secrets.Get method.
func (f *Fake) Get(ctx context.Context, name, namespace string) (*secrets.Secret, error) {
	args := f.Called(ctx, name, namespace)
	s := args.Get(0).(*secrets.Secret)
	return s, args.Error(1)
}

// NewFakeSecrets initializes a new fake implementation for secrets.
func NewFakeSecrets() *Fake {
	return &Fake{
		Mock: new(mock.Mock),
	}
}
