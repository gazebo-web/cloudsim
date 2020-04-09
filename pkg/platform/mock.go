package platform

type Mock struct {
	*Platform
}

// NewMock creates a mocked platform to be used by tests.
func NewMock(config Config) IPlatform {
	p := &Mock{}
	p.Config = config
	p.Platform = New(config)
	return p
}
