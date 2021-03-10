package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// SimulationServiceAdaptor implements the simulations.Service interface.
// It acts as an adaptor between the legacy code and the new code introduced in the code refactor.
type SimulationServiceAdaptor struct {
	db *gorm.DB
}

// UpdateScore updates a simulation's score.
func (sa *SimulationServiceAdaptor) UpdateScore(groupID simulations.GroupID, score *float64) error {
	panic("implement me")
}

// MarkStopped marks a simulation as stopped.
func (sa *SimulationServiceAdaptor) MarkStopped(groupID simulations.GroupID) error {
	panic("implement me")
}

// GetWebsocketToken returns a simulation's websocket authorization token.
func (sa *SimulationServiceAdaptor) GetWebsocketToken(groupID simulations.GroupID) (string, error) {
	panic("implement me")
}

// Get gets a simulation deployment with the given GroupID.
func (sa *SimulationServiceAdaptor) Get(groupID simulations.GroupID) (simulations.Simulation, error) {
	result, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetParent gets the parent simulation for the given GroupID.
func (sa *SimulationServiceAdaptor) GetParent(groupID simulations.GroupID) (simulations.Simulation, error) {
	gid := groupID.String()
	parent, err := GetParentSimulation(sa.db, &SimulationDeployment{GroupID: &gid})
	if err != nil {
		return nil, err
	}
	return parent, nil
}

// UpdateStatus persists the given status that assigns to the given GroupID.
func (sa *SimulationServiceAdaptor) UpdateStatus(groupID simulations.GroupID, status simulations.Status) error {
	dep, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return err
	}
	em := dep.updateSimDepStatus(sa.db, convertStatus(status))
	if em != nil {
		return em.BaseError
	}
	return nil
}

// Update updates the simulation identified by groupID with the data given in simulation.
func (sa *SimulationServiceAdaptor) Update(groupID simulations.GroupID, simulation simulations.Simulation) error {
	q := sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID.String()).Updates(simulation)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

// GetRobots returns the list of robots for the given groupID.
func (sa *SimulationServiceAdaptor) GetRobots(groupID simulations.GroupID) ([]simulations.Robot, error) {
	dep, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return nil, err
	}
	info, err := ReadExtraInfoSubT(dep)
	if err != nil {
		return nil, err
	}
	var result []simulations.Robot
	for _, robot := range info.Robots {
		r := new(SubTRobot)
		*r = robot
		result = append(result, r)
	}
	return result, nil
}

// NewSubTSimulationServiceAdaptor initializes a new simulations.Service implementation using the SimulationServiceAdaptor.
func NewSubTSimulationServiceAdaptor(db *gorm.DB) simulations.Service {
	return &SimulationServiceAdaptor{db: db}
}
