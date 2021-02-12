package simulations

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type SubTServiceAdaptor struct {
	db *gorm.DB
}

// Get gets a simulation deployment with the given GroupID.
func (sa *SubTServiceAdaptor) Get(groupID simulations.GroupID) (simulations.Simulation, error) {
	result, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetParent gets the parent simulation for the given GroupID.
func (sa *SubTServiceAdaptor) GetParent(groupID simulations.GroupID) (simulations.Simulation, error) {
	gid := groupID.String()
	parent, err := GetParentSimulation(sa.db, &SimulationDeployment{GroupID: &gid})
	if err != nil {
		return nil, err
	}
	return parent, nil
}

// UpdateStatus persists the given status that assigns to the given GroupID.
func (sa *SubTServiceAdaptor) UpdateStatus(groupID simulations.GroupID, status simulations.Status) error {
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

func (sa *SubTServiceAdaptor) Update(groupID simulations.GroupID, simulation simulations.Simulation) error {
	q := sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID.String()).Updates(simulation)
	if q.Error != nil {
		return q.Error
	}
	return nil
}

func (sa *SubTServiceAdaptor) GetRobots(groupID simulations.GroupID) ([]simulations.Robot, error) {
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

func NewSubTServiceAdaptor(db *gorm.DB) simulations.Service {
	return &SubTServiceAdaptor{db: db}
}
