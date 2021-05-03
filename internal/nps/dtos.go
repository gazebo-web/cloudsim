package nps

import (
	"time"
)

type GetSimulationResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	GroupID   string    `json:"groupid"`
	Status    string    `json:"status"`
	// The docker to run
	Image string `json:"image"`
	Args  string `json:"args"`
	URI   string `json:"uri"`
}

// StartRequest is the request used to start a simulation.
type StartRequest struct {
	// image is the docker image to run
	Image string `form:"image"`
	Args  string `form:"args"`
	Name  string `form:"name"`
}

// StartResponse is the response to the request of starting a simulation.
type StartResponse struct {
	Message    string `json:"message"`
	Simulation GetSimulationResponse
}

// StopRequest is the request to stop a simulation.
type StopRequest struct {
	GroupID string `json:"groupid"`
}

// StopResponse is the response to the request of stopping a simulation.
type StopResponse struct {
	Message    string `json:"message"`
	Simulation GetSimulationResponse
}

// ListRequest is the request to stop a simulation.
type ListRequest struct {
}

// ListResponse is the response to the request of stopping a simulation.
type ListResponse struct {
	Simulations []GetSimulationResponse
}

// AddUserRequest is the request used to add a registerd user.
type AddModifyUserRequest struct {
	Username        string `form:"username"`
	SimulationLimit int    `form:"simulation_limit"`
}

// AddUserResponse is the response used to add a registerd user.
type AddModifyUserResponse struct {
	Username        string `json:"username"`
	SimulationLimit int    `json:"simulation_limit"`
}
