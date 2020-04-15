package simulator

type IController interface {}

type Controller struct {
	Service IService
}

// NewController returns a new IController implementation.
func NewController(service IService) IController {
	var c IController
	c = &Controller{
		Service: service,
	}
	return c
}