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
	// Default value: false
	Public bool `json:"public"`
}

func (Track) TableName() string {
	return "subt_tracks"
}
