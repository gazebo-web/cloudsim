package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/groups"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
)

// IService
type IService interface {

}

// Service
type Service struct {
	nodeRepository nodes.IRepository
}

// NewSimulatorService
func NewSimulatorService(node nodes.IRepository) *Service {
	return &Service{
		nodeRepository:  node,
	}
}