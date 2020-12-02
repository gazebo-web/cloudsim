package secrets

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Fake struct {
	*mock.Mock
}

func (f *Fake) Get(ctx context.Context, name, namespace string) (*Secret, error) {
	args := f.Called(ctx, name, namespace)
	s := args.Get(0).(*Secret)
	return s, args.Error(1)
}

func NewFakeSecrets() *Fake {
	return &Fake{
		Mock: new(mock.Mock),
	}
}
