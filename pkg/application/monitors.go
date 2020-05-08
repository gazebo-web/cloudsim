package application

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// checkForExpiredSimulations is an internal helper that tests all the runningSimulations
// to check if they were alive more than expected, and in that case, schedules their termination.
func (app *Application) checkForExpiredSimulations() error {
	app.Platform().Simulator.RLock()
	defer app.Platform().Simulator.RUnlock()

	runningSims := app.Platform().Simulator.GetRunningSimulations()

	for groupID := range runningSims {
		rs := runningSims[groupID]

		if rs.IsExpired() || rs.Finished {
			app.Platform().RequestTermination(app.Platform().Context, groupID)
			reason := "expired"
			if rs.Finished {
				reason = "finished"
			}
			logger.Logger(app.Platform().Context).Info(fmt.Sprintf("Scheduled automatic termination of %s simulation: %s", reason, groupID))
		}
	}
	return nil
}

// updateMultiSimStatuses updates the the statuses of parent simulations from their children.
func (app *Application) updateMultiSimStatuses() error {
	parents, err := app.Services.Simulation.GetAllParentsWithErrors(
		simulations.StatusPending,
		simulations.StatusTerminatingInstances,
		[]simulations.ErrorStatus{simulations.ErrWhenInitializing, simulations.ErrWhenTerminating},
	)
	if err != nil {
		return err
	}
	for _, p := range *parents {
		if _, err := app.Services.Simulation.UpdateParentFromChildren(&p); err != nil {
			return err.BaseError
		}
	}
	return nil
}
