package metadata

import "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/robots"

type Metadata struct {
	Circuit string         `json:"circuit"`
	Robots  []robots.Robot `json:"robots,omitempty"`
}
