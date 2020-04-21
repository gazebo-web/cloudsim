package robots

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/robots"
)

// Robot represents a SubT Robot.
type Robot struct {
	gorm.Model
	robots.Robot
	Circuit string
}
