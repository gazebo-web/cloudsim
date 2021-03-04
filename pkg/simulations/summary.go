package simulations

import "github.com/jinzhu/gorm"

// Summary contains the total score and average statistics for a certain simulation or group of simulations.
type Summary struct {
	gorm.Model
	// GroupID identifies a simulation.
	GroupID *GroupID `json:"-" gorm:"not null;unique"`
	// Score is the simulation score.
	Score float64 `json:"-"`
	// SimTimeDurationAvg is the average of the simulation time duration.
	SimTimeDurationAvg float64 `json:"sim_time_duration_avg"`
	// SimTimeDurationStdDev is the standard deviation of the simulation time duration.
	SimTimeDurationStdDev  float64 `json:"sim_time_duration_std_dev" gorm:"-"`
	RealTimeDurationAvg    float64 `json:"real_time_duration_avg"`
	RealTimeDurationStdDev float64 `json:"real_time_duration_std_dev" gorm:"-"`
	ModelCountAvg          float64 `json:"model_count_avg"`
	ModelCountStdDev       float64 `json:"model_count_std_dev" gorm:"-"`
	Sources                string  `json:"-"`
}
