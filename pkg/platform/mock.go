package platform

import "context"

type Mock struct {
	StartMock              func(ctx context.Context) error
	StopMock               func(ctx context.Context) error
	RequestLaunchMock      func(ctx context.Context, groupID string)
	RequestTerminationMock func(ctx context.Context, groupID string)
}

func (m *Mock) Name() string {
	return "platform_test"
}

func (m *Mock) Start(ctx context.Context) error {
	return m.StartMock(ctx)
}

func (m *Mock) Stop(ctx context.Context) error {
	return m.StopMock(ctx)
}

func (m *Mock) RequestLaunch(ctx context.Context, groupID string) {
	m.RequestLaunchMock(ctx, groupID)
}

func (m *Mock) RequestTermination(ctx context.Context, groupID string) {
	m.RequestTerminationMock(ctx, groupID)
}

// NewMock creates a mocked platform to be used by tests.
func NewMock() *Mock {
	p := &Mock{}
	return p
}
