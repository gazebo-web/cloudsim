package nps

import (
	"time"
)

type GetSimulationResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	GroupID   string    `json:"groupid"`
	Status    string    `json:"status"`
	// The docker to run
	Image string `json:"image"`
	Args  string `json:"args"`
	URI   string `json:"uri"`
	IP    string `json:"ip"`
}

// StartRequest is the request used to start a simulation.
type StartRequest struct {
	// image is the docker image to run
	Image string `form:"image"`
	Args  string `form:"args"`
}

// StartResponse is the response to the request of starting a simulation.
type StartResponse struct {
	Message    string `json:"message"`
	Simulation GetSimulationResponse
}

// StopRequest is the request to stop a simulation.
type StopRequest struct {
}

// StopResponse is the response to the request of stopping a simulation.
type StopResponse struct {
}

// ListRequest is the request to stop a simulation.
type ListRequest struct {
}

// ListResponse is the response to the request of stopping a simulation.
type ListResponse struct {
	Simulations []GetSimulationResponse
}
