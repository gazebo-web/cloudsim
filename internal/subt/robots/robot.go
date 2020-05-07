package robots

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/robots"
)

// Robot represents a SubT Robot.
type Robot struct {
	robots.Robot
	Credits   int    `json:"credits"`
	Thumbnail string `json:"thumbnail"`
}

type Robots []Robot

func (Robot) TableName() string {
	return "subt_robots"
}
