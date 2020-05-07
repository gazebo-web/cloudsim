package robots

// Robot represents a generic robot to be used on simulations.
type Robot struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Type  string `json:"type"`
}
