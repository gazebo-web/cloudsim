package simulator

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// createCheckPendingStatusJob returns an job that will be used to check if a simulation is pending or not.
func (s *subTSimulator) createCheckPendingStatusJob(groupID simulations.GroupID) *actions.Job {
	job := actions.JobFunc(func(tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		gid, ok := value.(simulations.GroupID)
		if !ok {
			return nil, simulations.ErrInvalidGroupID
		}
		sim, err := s.simulationService.Get(gid)
		if err != nil {
			return nil, err
		}
		if sim.Status() != simulations.StatusPending {
			return nil, simulations.ErrIncorrectStatus
		}
		return gid, nil
	})

	return &actions.Job{
		Name:            fmt.Sprintf("check-pending-status-%s", groupID),
		PreHooks:        nil,
		Execute:         job,
		PostHooks:       nil,
		RollbackHandler: nil,
		InputType:       actions.GetJobDataType(groupID),
		OutputType:      actions.GetJobDataType(groupID),
	}
}

// createCheckSimulationIsParentJob returns a job that will be used to check if a simulation is not a parent.
func (s *subTSimulator) createCheckSimulationIsParentJob(groupID simulations.GroupID) *actions.Job {
	job := actions.JobFunc(func(tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		gid, ok := value.(simulations.GroupID)
		if !ok {
			return nil, simulations.ErrInvalidGroupID
		}
		sim, err := s.simulationService.Get(gid)
		if err != nil {
			return nil, err
		}
		if sim.Kind() == simulations.SimParent {
			_, err := s.simulationService.Reject(gid)
			if err != nil {
				return nil, err
			}
			return nil, simulations.ErrIncorrectKind
		}
		return gid, nil
	})

	return &actions.Job{
		Name:            fmt.Sprintf("check-simulation-parenthood-%s", groupID),
		PreHooks:        nil,
		Execute:         job,
		PostHooks:       nil,
		RollbackHandler: nil,
		InputType:       actions.GetJobDataType(groupID),
		OutputType:      actions.GetJobDataType(groupID),
	}
}

// createCheckParentSimulationWithErrorJob returns a job that will be used to check if a parent simulation
// from the given children simulation has an error or not.
func (s *subTSimulator) createCheckParentSimulationWithErrorJob(groupID simulations.GroupID) *actions.Job {
	job := actions.JobFunc(func(tx *gorm.DB, deployment *actions.Deployment, value interface{}) (interface{}, error) {
		gid, ok := value.(simulations.GroupID)
		if !ok {
			return nil, simulations.ErrInvalidGroupID
		}
		sim, err := s.simulationService.Get(gid)
		if err != nil {
			return nil, err
		}
		if sim.Kind() != simulations.SimChild {
			return gid, nil
		}
		parent, err := s.simulationService.GetParent(gid)
		if err != nil {
			return nil, err
		}
		if simerr := parent.Error(); simerr != nil {
			return nil, simulations.ErrParentSimulationWithError
		}
		return gid, nil
	})

	return &actions.Job{
		Name:            fmt.Sprintf("check-parent-simulation-with-error-%s", groupID),
		PreHooks:        nil,
		Execute:         job,
		PostHooks:       nil,
		RollbackHandler: nil,
		InputType:       actions.GetJobDataType(groupID),
		OutputType:      actions.GetJobDataType(groupID),
	}
}
