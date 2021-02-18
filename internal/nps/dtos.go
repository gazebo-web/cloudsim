package nps

// StartRequest is the request used to start a simulation.
type StartRequest struct {
	// image is the docker image to run
	Image string `form:"image"`
}

// StartResponse is the response to the request of starting a simulation.
type StartResponse struct {
}

// StopRequest is the request to stop a simulation.
type StopRequest struct {
}

// StopResponse is the response to the request of stopping a simulation.
type StopResponse struct {
}
