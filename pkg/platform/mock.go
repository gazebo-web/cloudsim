package platform

type PlatformMock struct {
	Platform
}

func NewMock(config Config) PlatformMock {
	p := PlatformMock{}
	p.Config = config
	p.Platform = New(config)
	return p
}
