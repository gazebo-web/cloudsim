package platform

type Mock struct {
	Platform
}

func NewMock(config Config) IPlatform {
	p := &Mock{}
	p.Config = config
	p.Platform = New(config)
	return p
}
