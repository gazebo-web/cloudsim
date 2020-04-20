package stats

// Statistics contains the summary values of a simulation run.
type Statistics struct {
	WasStarted          int `yaml:"was_started"`
	SimTimeDurationSec  int `yaml:"sim_time_duration_sec"`
	RealTimeDurationSec int `yaml:"real_time_duration_sec"`
	ModelCount          int `yaml:"model_count"`
}
