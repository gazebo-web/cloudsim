package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/groups"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
)

type IService interface {

}

type Service struct {
	groupRepository groups.IRepository
	nodeRepository nodes.IRepository
}