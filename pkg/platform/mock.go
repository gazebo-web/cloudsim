package platform

type PlatformMock struct {
	Platform
}

func NewMock(config Config) IPlatform {
	p := &PlatformMock{}
	p.Config = config
	p.Platform = New(config)
	return p
}
