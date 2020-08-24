package simulations

import (
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type GetAllInput struct {
	p               *ign.PaginationRequest
	byStatus        *Status
	invertStatus    bool
	byErrStatus     *ErrorStatus
	invertErrStatus bool
	user            *fuel.User
	includeChildren bool
}
