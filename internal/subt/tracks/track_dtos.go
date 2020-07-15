package tracks

// CreateTrack is a Data Access Object for creating a new track.
type CreateTrack struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	BridgeImage string `json:"bridge_image"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string `json:"stats_topic"`
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic string `json:"warmup_topic"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int `json:"max_sim_seconds"`
	// Public sets a track open for single runs.
	// Default value: false
	Public bool `json:"public"`
}

// UpdateTrack is a Data Access Object for updating an existent track.
type UpdateTrack struct {
	Name        string
	Image       string
	BridgeImage string
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic string
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int
	// Public sets a track open for single runs.
	// Default value: false
	Public bool
}
