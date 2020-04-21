package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/groups"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
)

type IService interface {

}

type Service struct {
	nodeRepository nodes.IRepository
}

func NewSimulatorService(node nodes.IRepository) *Service {
	return &Service{
		nodeRepository:  node,
	}
}