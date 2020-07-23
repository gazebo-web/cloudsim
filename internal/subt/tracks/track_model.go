package tracks

import (
	"github.com/jinzhu/gorm"
)

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
}

// SingularName returns the Track entity name in singular.
func (Track) SingularName() string {
	return "Track"
}

// PluralName returns the Track entity name in plural.
func (Track) PluralName() string {
	return "Tracks"
}

// TableName returns the name of the Track table.
func (Track) TableName() string {
	return "subt_tracks"
}

// CreateTrackFromInput receives an input and returns a new Track with the input values.
func CreateTrackFromInput(input CreateTrackInput) *Track {
	return &Track{
		Name:          input.Name,
		Image:         input.Image,
		BridgeImage:   input.BridgeImage,
		StatsTopic:    input.StatsTopic,
		WarmupTopic:   input.WarmupTopic,
		MaxSimSeconds: input.MaxSimSeconds,
		Public:        input.Public,
	}
}
