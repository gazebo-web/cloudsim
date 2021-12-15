package simulations

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// SimulationServiceAdaptor implements the simulations.Service interface.
// It acts as an adaptor between the legacy code and the new code introduced in the code refactor.
type SimulationServiceAdaptor struct {
	db *gorm.DB
}

// UpdateScore updates the score of a certain simulation deployment.
func (sa *SimulationServiceAdaptor) UpdateScore(groupID simulations.GroupID, score *float64) error {
	dep, err := GetSimulationDeployment(sa.db, groupID.String())
	if err != nil {
		return err
	}

	em := dep.UpdateScore(sa.db, score)
	if em != nil {
		return em.BaseError
	}

	return nil
}

// MarkStopped marks a simulation with the time when it has stopped running.
// If the StoppedAt value is already set, it won't be updated.
func (sa *SimulationServiceAdaptor) MarkStopped(groupID simulations.GroupID) error {
	return sa.db.Model(&SimulationDeployment{}).Transaction(func(tx *gorm.DB) error {
		var sim SimulationDeployment
		if err := sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID).First(&sim).Error; err != nil {
			return err
		}
		if sim.StoppedAt != nil {
			return nil
		}
		at := time.Now()
		if err := sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID).Update(SimulationDeployment{
			StoppedAt: &at,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

// Create creates a simulation (SimulationDeployment) from the given input.
func (sa *SimulationServiceAdaptor) Create(input simulations.CreateSimulationInput) (simulations.Simulation, error) {
	dep, err := NewSimulationDeployment()
	if err != nil {
		return nil, err
	}

	gid := simulations.GroupID(uuid.NewV4().String()).String()

	image := SliceToStr(input.Image)

	dep.Owner = input.Owner
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

// GetWebsocketToken returns a simulation's websocket authorization token.
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
	return GetSimulationDeployment(sa.db, groupID.String())
}

// GetParent gets the parent simulation for the given GroupID.
func (sa *SimulationServiceAdaptor) GetParent(groupID simulations.GroupID) (simulations.Simulation, error) {
	gid := groupID.String()
	return GetParentSimulation(sa.db, &SimulationDeployment{ParentGroupID: &gid})
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
	return sa.db.Model(&SimulationDeployment{}).Where("group_id = ?", groupID.String()).Updates(simulation).Error
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
