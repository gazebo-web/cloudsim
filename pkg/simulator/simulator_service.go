package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
)

// IService
type IService interface {
}

// Service
type Service struct {
	nodeRepository nodes.Repository
}

// NewSimulatorService
func NewSimulatorService(node nodes.Repository) *Service {
	return &Service{
		nodeRepository: node,
	}
}
