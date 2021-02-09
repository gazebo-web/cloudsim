package simulations

import "github.com/jinzhu/gorm"

// Summary contains the total score and average statistics for a certain simulation or group of simulations.
type Summary struct {
	gorm.Model
	// GroupID identifies a simulation.
	GroupID *GroupID `json:"-" gorm:"not null;unique"`
	// Score is the simulation score.
	Score *float64 `json:"-"`
	// SimTimeDurationAvg is the average value of the simulation time duration.
	SimTimeDurationAvg float64 `json:"sim_time_duration_avg"`
	// SimTimeDurationStdDev is the standard deviation value of the simulation time duration. Only used by simulations.SimParent.
	SimTimeDurationStdDev float64 `json:"sim_time_duration_std_dev" gorm:"-"`
	// RealTimeDurationAvg is the average value of the real time duration.
	RealTimeDurationAvg float64 `json:"real_time_duration_avg"`
	// RealTimeDurationStdDev is the standard deviation value of the real time duration. Only used by simulations.SimParent.
	RealTimeDurationStdDev float64 `json:"real_time_duration_std_dev" gorm:"-"`
	// ModelCountAvg is the average value of the model count.
	ModelCountAvg float64 `json:"model_count_avg"`
	// ModelCountStdDev is the standard deviation value of the model count. Only used by simulations.SimParent.
	ModelCountStdDev float64 `json:"model_count_std_dev" gorm:"-"`
	// Sources is used save all the children group ids. Only used by simulations.SimParent.
	Sources string `json:"-"`
	// RunData has the simulation run data.
	RunData string `json:"-" gorm:"-"`
}
