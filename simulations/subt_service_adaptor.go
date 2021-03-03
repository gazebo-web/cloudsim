package simulations

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// SimulationServiceAdaptor implements the simulations.Service interface.
// It acts as an adaptor between the legacy code and the new code introduced in the code refactor.
type SimulationServiceAdaptor struct {
	db *gorm.DB
}

// MarkStopped marks a simulation with the time where it has stopped running.
func (sa *SimulationServiceAdaptor) MarkStopped(groupID simulations.GroupID) error {
	at := time.Now()
	if err := sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID).Update(SimulationDeployment{
		StoppedAt: &at,
	}).Error; err != nil {
		return err
	}
	return nil
}

// Create creates a simulation (SimulationDeployment) from the given input.
func (sa *SimulationServiceAdaptor) Create(input simulations.CreateSimulationInput) (simulations.Simulation, error) {
	dep, err := NewSimulationDeployment()
	if err != nil {
		return nil, err
	}

	gid := simulations.GroupID(uuid.NewV4().String()).String()

	image := SliceToStr(input.Image)

	dep.Owner = &input.Owner
	dep.Name = &input.Name
	dep.Creator = &input.Creator
	dep.Private = &input.Private
	dep.StopOnEnd = &input.StopOnEnd
	dep.Image = &image
	dep.GroupID = &gid
	dep.DeploymentStatus = simPending.ToPtr()
	dep.Extra = &input.Extra
	dep.ExtraSelector = &input.Track
	dep.Robots = &input.Robots
	dep.Held = false

	if err := sa.db.Model(&SimulationDeployment{}).Create(dep).Error; err != nil {
		return nil, err
	}
	return dep, nil
}

// GetWebsocketToken gets the authorization token to connect a websocket server for the given simulation.
func (sa *SimulationServiceAdaptor) GetWebsocketToken(groupID simulations.GroupID) (string, error) {
	dep, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return "", err
	}
	if dep.AuthorizationToken == nil {
		return "", errors.New("missing access token")
	}
	return *dep.AuthorizationToken, nil
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
