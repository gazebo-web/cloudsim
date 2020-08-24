package robots

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/robots"
)

// Robot represents a SubT Robot.
type Robot struct {
	gorm.Model
	robots.Robot
	Credits   int    `json:"credits"`
	Thumbnail string `json:"thumbnail"`
}

type Robots []Robot

func (Robot) TableName() string {
	return "subt_robots"
}

type RobotConfig struct {
	gorm.Model
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Type      string `json:"type"`
	Credits   int    `json:"credits"`
	Thumbnail string `json:"thumbnail"`
}

func (RobotConfig) TableName() string {
	return "subt_robot_configs"
}
