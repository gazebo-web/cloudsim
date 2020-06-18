package simulator

// Controller
type Controller interface{}

// controller
type controller struct {
	Service IService
}

// NewController returns a new Controller implementation.
func NewController(service IService) Controller {
	var c Controller
	c = &controller{
		Service: service,
	}
	return c
}
