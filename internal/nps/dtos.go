package nps

// StartRequest is the request used to start a simulation.
type StartRequest struct {
	// image is the docker image to run
	Image string `form:"image"`
  Args string `form:"args"`
}

// StartResponse is the response to the request of starting a simulation.
type StartResponse struct {
  URI string `json:"uri"`
}

// StopRequest is the request to stop a simulation.
type StopRequest struct {
}

// StopResponse is the response to the request of stopping a simulation.
type StopResponse struct {
}
