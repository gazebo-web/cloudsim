package tracks

import "github.com/jinzhu/gorm"

// Track is a world that will be used to run a simulation.
type Track struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
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
	// Seed is the world's seed of this track. It no seed is provided, a random one will be used instead.
	Seed *int `json:"seed"`
	// World is the world
	World string
}

// TableName returns the name of the Track table.
func (Track) TableName() string {
	return "subt_tracks"
}

// CreateTrackFromInput receives an input and returns a new Track with the input values.
func CreateTrackFromInput(input CreateTrackInput) Track {
	return Track{
		Name:          input.Name,
		Image:         input.Image,
		BridgeImage:   input.BridgeImage,
		StatsTopic:    input.StatsTopic,
		WarmupTopic:   input.WarmupTopic,
		MaxSimSeconds: input.MaxSimSeconds,
		Public:        input.Public,
		World:         input.World,
	}
}

// UpdateTrackFromInput receives a model and an updated input.
// It returns the model updated with the input values.
func UpdateTrackFromInput(model Track, input UpdateTrackInput) Track {
	model.Name = input.Name
	model.Image = input.Image
	model.BridgeImage = input.BridgeImage
	model.StatsTopic = input.StatsTopic
	model.WarmupTopic = input.WarmupTopic
	model.MaxSimSeconds = input.MaxSimSeconds
	model.Public = input.Public
	model.World = input.World
	return model
}
