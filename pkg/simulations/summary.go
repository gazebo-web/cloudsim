package simulations

type Summary struct {
	Score                  float64 `json:"-"`
	SimTimeDurationAvg     float64 `json:"sim_time_duration_avg"`
	SimTimeDurationStdDev  float64 `json:"sim_time_duration_std_dev"`
	RealTimeDurationAvg    float64 `json:"real_time_duration_avg"`
	RealTimeDurationStdDev float64 `json:"real_time_duration_std_dev"`
	ModelCountAvg          float64 `json:"model_count_avg"`
	ModelCountStdDev       float64 `json:"model_count_std_dev"`
	Sources                string  `json:"-"`
}
