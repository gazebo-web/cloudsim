package tracks

// CreateTrackInput is an input for creating a new track.
type CreateTrackInput struct {
	Name        string `json:"name" validate:"required,gt=10"`
	Image       string `json:"image" validate:"required"`
	BridgeImage string `json:"bridge_image" validate:"required"`
	// Topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string `json:"stats_topic" validate:"required"`
	// Topic used to track when the simulation officially starts and ends
	WarmupTopic string `json:"warmup_topic" validate:"required"`
	// Maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int `json:"max_sim_seconds" validate:"required"`
	// Public makes a track available for launching directly.
	// Tracks that are not public can only be launched as part of a Circuit.
	Public bool `json:"public" validate:"required"`
}

// UpdateTrackInput is an input for updating an existent track.
type UpdateTrackInput struct {
	Name        string `json:"name" validate:"required,gt=10"`
	Image       string `json:"image" validate:"required"`
	BridgeImage string `json:"bridge_image" validate:"required"`
	// StatsTopic is a topic used to track general stats of the simulation (runtime, sim runtime, etc.)
	StatsTopic string `json:"stats_topic" validate:"required"`
	// WarmupTopic is a topic used to track when the simulation officially starts and ends
	WarmupTopic string `json:"warmup_topic" validate:"required"`
	// MaxSimSeconds is the maximum number of allowed "simulation seconds" for each world. 0 means unlimited.
	MaxSimSeconds int `json:"max_sim_seconds" validate:"required"`
	// Public makes a track available for launching directly.
	// Tracks that are not public can only be launched as part of a Circuit.
	Public bool `json:"public" validate:"required"`
}
