package simulations

// Statistics contains the summary values of a simulation run.
type Statistics struct {
	// Started is true if the simulation was started.
	Started int `yaml:"was_started"`
	// SimulationTime is the duration in seconds of the simulation time.
	SimulationTime int `yaml:"sim_time_duration_sec"`
	// RealTime is the real duration in seconds of the simulation.
	RealTime int `yaml:"real_time_duration_sec"`
	// ModelCount is the amount of models used in the simulation.
	ModelCount int `yaml:"model_count"`
}
