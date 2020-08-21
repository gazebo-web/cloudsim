package application

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// Services groups the services needed by an application to launch simulations.
type Services interface {
	// Simulations returns a service to operate over different simulations.
	Simulations() simulations.Service
}
