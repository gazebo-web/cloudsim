package metadata

import (
	"encoding/json"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/robots"
)

type Metadata struct {
	Circuit string         `json:"circuit"`
	Robots  []robots.Robot `json:"robots,omitempty"`
	RunIndex *int `json:"run_index,omitempty"`
	WorldIndex *int        `json:"world_index,omitempty"`
}

func (m Metadata) ToJSON() (*string, error) {
	result := new(string)
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	*result = string(b)
	return result, nil
}