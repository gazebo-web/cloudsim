package tracks

// CreateTrackInput is an input for creating a new track.
type CreateTrackInput struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	BridgeImage string `json:"bridge_image"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string `json:"stats_topic"`
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic string `json:"warmup_topic"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int `json:"max_sim_seconds"`
	// Public makes a track available for launching directly.
	// Tracks that are not public can only be launched as part of a Circuit.
	Public bool `json:"public"`
}

// UpdateTrackInput is an input for updating an existent track.
type UpdateTrackInput struct {
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
	Public bool `json:"public"`
}
