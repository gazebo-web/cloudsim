package simulations

import (
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	Start()
	LaunchHeld()
	Restart()
	Shutdown()
	GetAll()
	Get()
	GetDownloadableLogs()
	GetLiveLogs()
}

// Controller is an IController implementation.
type Controller struct {
	services services
	formDecoder *form.Decoder
	validator *validator.Validate
}

type services struct {
	Simulation IService
	User users.IService
}

func (c *Controller) Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	var createSim simulations.SimulationCreate
	if em := tools.ParseFormStruct(&createSim, r, c.formDecoder); em != nil {
		return nil, em
	}

	if em := tools.ValidateStruct(&createSim, c.validator); em != nil {
		return nil, em
	}

	// HACK until we have more platforms and applications.
	// TODO: remove this
	if createSim.Platform == "" {
		createSim.Platform = getDefaultPlatformName()
	}
	if createSim.Application == "" {
		createSim.Application = getDefaultApplicationName()
	}

	// Set the owner, if missing
	if createSim.Owner == "" {
		createSim.Owner = *user.Username
	}

	// Allow the custom Application to customize the SimulationCreate request
	if em := SimServImpl.CustomizeSimRequest(r.Context(), r, tx, &createSim, *user.Username); em != nil {
		return nil, em
	}

	return SimServImpl.StartSimulationAsync(r.Context(), tx, &createSim, user)
}

func (c *Controller) LunchHeld() {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.LaunchSimulationAsync(r.Context(), tx, groupID, user)
}

func (c *Controller) Restart() {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.RestartSimulationAsync(r.Context(), tx, groupID, user)
}

func (c *Controller) Shutdown() {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.ShutdownSimulationAsync(r.Context(), tx, groupID, user)
}

func (c *Controller) GetAll() {
	// Prepare pagination
	pr, em := ign.NewPaginationRequest(r)
	if em != nil {
		return nil, em
	}

	// Get the parameters
	params := r.URL.Query()
	var status *DeploymentStatus
	invertStatus := false
	invertErrStatus := false
	if len(params["status"]) > 0 && len(params["status"][0]) > 0 {
		statusStr := params["status"][0]
		invertStatus = statusStr[0] == '!'
		sliceIndex := 0
		if invertStatus {
			sliceIndex = 1
		}
		status = DeploymentStatusFrom(statusStr[sliceIndex:])
	}
	var errStatus *ErrorStatus
	if len(params["errorStatus"]) > 0 && len(params["errorStatus"][0]) > 0 {
		statusStr := params["errorStatus"][0]
		invertErrStatus = statusStr[0] == '!'
		sliceIndex := 0
		if invertErrStatus {
			sliceIndex = 1
		}
		errStatus = ErrorStatusFrom(statusStr[sliceIndex:])
	}

	includeChildren := false
	if len(params["children"]) > 0 && len(params["children"][0]) > 0 {
		if flag, err := strconv.ParseBool(params["children"][0]); err == nil {
			includeChildren = flag
		}
	}

	// TODO: This is SubT specific and should be moved
	var circuit *string
	if len(params["circuit"]) > 0 && len(params["circuit"][0]) > 0 {
		circuit = &params["circuit"][0]
	}

	sims, pagination, em := SimServImpl.SimulationDeploymentList(
		r.Context(),
		pr,
		tx,
		status,
		invertStatus,
		errStatus,
		invertErrStatus,
		circuit,
		user,
		sptr(getDefaultApplicationName()),
		includeChildren,
	)
	if em != nil {
		return nil, em
	}

	ign.WritePaginationHeaders(*pagination, w, r)
}

func (c *Controller) Get() {

}

func GetDownloadableLogs() {

}

func GetLiveLogs() {

}