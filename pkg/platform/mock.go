package platform

type Mock struct {
	IPlatform
	IPlatformSetup
}

// NewMock creates a mocked platform to be used by tests.
func NewMock(config Config) IPlatform {
	var p IPlatform
	p = &Mock{}
	p = New(config)
	return p
}

// NewSetupMock creates a mocked platform initializer to be used by tests.
func NewSetupMock(config Config) IPlatformSetup {
	var p IPlatformSetup
	p = &Mock{}
	p = New(config)
	return p
}
