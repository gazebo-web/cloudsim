package platform

import "context"

type IPlatformService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func (p Platform) Start(ctx context.Context) error {
	return nil
}

func (p Platform) Stop(ctx context.Context) error {
	return nil
}