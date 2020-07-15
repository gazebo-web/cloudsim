package tracks

import "github.com/jinzhu/gorm"

// Track is a world that will be used to run a simulation.
type Track struct {
	gorm.Model
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
