package metadata

import (
	"encoding/json"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/robots"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
)

type Metadata struct {
	Circuit    string         `json:"circuit"`
	Robots     []robots.Robot `json:"robots,omitempty"`
	RunIndex   *int           `json:"run_index,omitempty"`
	WorldIndex *int           `json:"world_index,omitempty"`
}

func (m Metadata) ToJSON() (*string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return tools.Sptr(string(b)), nil
}

func Read(sim *simulations.Simulation) (*Metadata, error) {
	var extra Metadata
	if err := json.Unmarshal([]byte(*sim.Extra), &extra); err != nil {
		return nil, err
	}
	return &extra, nil
}
