package simulator

// GazeboConfig represents a set of parameters that an application passes to the Gazebo Server.
type GazeboConfig struct {
	WorldStatsTopic string
	WorldWarmupTopic string
	MaxSeconds int64
}
