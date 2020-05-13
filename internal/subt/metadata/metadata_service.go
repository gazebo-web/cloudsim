package metadata

import (
	"encoding/json"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

func Read(sim *simulations.Simulation) (*Metadata, error) {
	var extra Metadata
	if err := json.Unmarshal([]byte(*sim.Extra), &extra); err != nil {
		return nil, err
	}
	return &extra, nil
}