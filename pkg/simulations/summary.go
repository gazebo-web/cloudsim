package simulations

import "github.com/jinzhu/gorm"

// Summary contains the total score and average statistics for a certain simulation or group of simulations
type Summary struct {
	gorm.Model
	GroupID                *GroupID `json:"-" gorm:"not null;unique"`
	Score                  float64  `json:"-"`
	SimTimeDurationAvg     float64  `json:"sim_time_duration_avg"`
	SimTimeDurationStdDev  float64  `json:"sim_time_duration_std_dev" gorm:"-"`
	RealTimeDurationAvg    float64  `json:"real_time_duration_avg"`
	RealTimeDurationStdDev float64  `json:"real_time_duration_std_dev" gorm:"-"`
	ModelCountAvg          float64  `json:"model_count_avg"`
	ModelCountStdDev       float64  `json:"model_count_std_dev" gorm:"-"`
	Sources                string   `json:"-"`
}
