package simulations

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"strconv"
)

// HTTPHandler is used to invoke inner logic based on incoming Http requests.
type HTTPHandler struct {
	UserAccessor useracc.Service
}

// HTTPHandlerInstance is the default HTTPHandler instance. It is used by routes.go.
var HTTPHandlerInstance *HTTPHandler

// NewHTTPHandler creates a new HTTPHandler.
func NewHTTPHandler(ctx context.Context, ua useracc.Service) (*HTTPHandler, error) {
	return &HTTPHandler{
		UserAccessor: ua,
	}, nil
}

type handlerWithUser func(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)

// WithUser is a middleware that checks for a valid user from the JWT and passes
// the user to the handlerWithUser.
// If the user is not present or invalid, an error is returned.
func WithUser(handler handlerWithUser) ign.HandlerWithResult {
	return func(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
		// Get JWT user. Fail if invalid or missing
		user, ok, em := HTTPHandlerInstance.UserAccessor.UserFromJWT(r)
		if !ok {
			return nil, em
		}
		return handler(user, tx, w, r)
	}
}

// WithUserOrAnonymous is a middleware that checks for a valid user from the JWT and passes
// the user to the handlerWithUser.
// If the user is not present or invalid, a nil user is passed.
func WithUserOrAnonymous(handler handlerWithUser) ign.HandlerWithResult {
	return func(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
		// Get JWT user. Fail if invalid or missing
		user, ok, _ := HTTPHandlerInstance.UserAccessor.UserFromJWT(r)
		if !ok {
			return handler(nil, tx, w, r)
		}
		return handler(user, tx, w, r)
	}
}

// getDefaultPlatformName returns the default platform for which to run
// simulations
func getDefaultApplicationName() string {
	// HACK This should be changed once more applications become available
	return applicationSubT
}

// getDefaultPlatformName returns the default platform on which to run
// simulations
func getDefaultPlatformName() string {
	// HACK This should be changed once more platforms become available
	return platformSubT
}

// CloudsimSimulationCreate is the main func to launch a new simulation
// You can request this method with the following cURL request:
//   curl -k -X POST --url http://localhost:8001/1.0/simulations
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudsimSimulationCreate(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	// Parse form's values and files.
	// https://golang.org/pkg/net/http/#Request.ParseMultipartForm
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	// CreateSimulation is the input form
	var createSim CreateSimulation
	if em := ParseStruct(&createSim, r, true); em != nil {
		return nil, em
	}

	// HACK until we have more platforms and applications.
	// TODO: remove this
	if createSim.Application == "" {
		createSim.Application = getDefaultApplicationName()
	}

	// Set the owner, if missing
	if createSim.Owner == "" {
		createSim.Owner = *user.Username
	}

	// Allow the custom Application to customize the CreateSimulation request
	if em := SimServImpl.CustomizeSimRequest(r.Context(), r, tx, &createSim, *user.Username); em != nil {
		return nil, em
	}

	return SimServImpl.StartSimulationAsync(r.Context(), tx, &createSim, user)
}

// CloudsimSimulationDelete finishes all resources associated to a cloudsim simulation.
// (eg. Nodes, Hosts, Pods)
// You can request this method with the following cURL request:
//   curl -k -X DELETE --url http://localhost:8001/1.0/simulations/{group}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudsimSimulationDelete(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.ShutdownSimulationAsync(r.Context(), tx, groupID, user)
}

// CloudsimSimulationLaunch launches a held simulation.
// You can request this method with the following cURL request:
//   curl -k -X POST --url http://localhost:8001/1.0/simulations/{group}/launch
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudsimSimulationLaunch(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.LaunchSimulationAsync(r.Context(), tx, groupID, user)
}

// CloudsimSimulationRestart restarts a failed single simulation.
// You can request this method with the following cURL request:
//   curl -k -X POST --url http://localhost:8001/1.0/simulations/{group}/restart
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudsimSimulationRestart(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.RestartSimulationAsync(r.Context(), tx, groupID, user)
}

// CloudsimSimulationList returns a list with simulation deployments.
// You can request this method with the following cURL request:
//   curl -k -X GET --url http://localhost:8001/1.0/simulations
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudsimSimulationList(user *users.User, tx *gorm.DB,
	w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

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

	var owner *string
	if len(params["owner"]) > 0 && len(params["owner"][0]) > 0 {
		owner = &params["owner"][0]
	}

	var private *bool
	if len(params["private"]) > 0 && len(params["private"][0]) > 0 {
		if flag, err := strconv.ParseBool(params["private"][0]); err == nil {
			private = &flag
		}
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
		owner,
		private,
	)
	if em != nil {
		return nil, em
	}

	ign.WritePaginationHeaders(*pagination, w, r)

	return sims, nil
}

// CustomRuleList returns a paginated list of circuit custom rules,
// filtering by circuit, owner or rule. This operation can only be performed by
// a system or application administrator.
// GET parameters include: application, circuit, owner and rule_type.
// You can request this method with the following cURL request:
//   curl -k -X GET --url http://localhost:8001/1.0/rules
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CustomRuleList(user *users.User, tx *gorm.DB,
	w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Prepare pagination
	pr, em := ign.NewPaginationRequest(r)
	if em != nil {
		return nil, em
	}

	// Get filter parameters
	params := r.URL.Query()
	var application *string
	if len(params["application"]) > 0 && len(params["application"][0]) > 0 {
		application = &params["application"][0]
	}
	var circuit *string
	if len(params["circuit"]) > 0 && len(params["circuit"][0]) > 0 {
		circuit = &params["circuit"][0]
	}
	var owner *string
	if len(params["owner"]) > 0 && len(params["owner"][0]) > 0 {
		owner = &params["owner"][0]
	}
	var ruleType *CustomRuleType
	if len(params["rule_type"]) > 0 && len(params["rule_type"][0]) > 0 {
		rule := CustomRuleType(params["rule_type"][0])
		ruleType = &rule
	}

	rules, pagination, em := SimServImpl.CustomRuleList(r.Context(), pr, tx, user, application, circuit, owner,
		ruleType)
	if em != nil {
		return nil, em
	}

	_ = ign.WritePaginationHeaders(*pagination, w, r)

	return rules, nil
}

// RemainingSubmissions contains GetRemaingSubmissions response struct
type RemainingSubmissions struct {
	RemainingSubmissions *int `json:"remaining_submissions"`
}

// GetRemainingSubmissions returns the number of remaining submissions for an
// owner in a circuit.
// You can request this method with the following cURL request:
//   curl -k -X GET --url http://localhost:8000/1.0/{circuit}/remaining_submissions/{owner}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func GetRemainingSubmissions(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	requestVars := mux.Vars(r)
	circuit, err := requestVars["circuit"]
	if !err {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	owner, err := requestVars["owner"]
	if !err {
		return nil, ign.NewErrorMessage(ign.ErrorOwnerNotInRequest)
	}

	remaining, em := SimServImpl.GetRemainingSubmissions(r.Context(), tx, user, &circuit, &owner)
	if em != nil {
		return nil, em
	}

	return RemainingSubmissions{remaining.(*int)}, nil
}

// SetCustomRule creates or updates a custom rule for an owner in a circuit.
// You can request this method with the following cURL request:
//   curl -k -X PUT --url http://localhost:8000/1.0/rules/{circuit}/{owner}/{rule}/{value}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func SetCustomRule(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	requestVars := mux.Vars(r)
	circuit, ok := requestVars["circuit"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	owner, ok := requestVars["owner"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorOwnerNotInRequest)
	}
	ruleString, ok := requestVars["rule"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	ruleType := CustomRuleType(ruleString)
	value, ok := requestVars["value"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	rule, em := SimServImpl.SetCustomRule(r.Context(), tx, user, sptr(getDefaultApplicationName()), &circuit, &owner,
		&ruleType, &value)

	return rule, em
}

// DeleteCustomRule deletes a custom rule for an owner in a circuit.
// You can request this method with the following cURL request:
//   curl -k -X DELETE --url http://localhost:8000/1.0/rules/{circuit}/{owner}/{rule}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func DeleteCustomRule(user *users.User, tx *gorm.DB, w http.ResponseWriter,
	r *http.Request) (interface{}, *ign.ErrMsg) {
	requestVars := mux.Vars(r)
	circuit, ok := requestVars["circuit"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	owner, ok := requestVars["owner"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorOwnerNotInRequest)
	}
	ruleString, ok := requestVars["rule"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	ruleType := CustomRuleType(ruleString)

	rule, em := SimServImpl.DeleteCustomRule(r.Context(), tx, user, sptr(getDefaultApplicationName()), &circuit,
		&owner, &ruleType)
	if em != nil {
		return nil, em
	}

	return rule, em
}

// GetCloudsimSimulation returns a single simulation.
// You can request this method with the following cURL request:
//   curl -k -X GET --url http://localhost:8000/1.0/simulations/{group}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func GetCloudsimSimulation(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	return SimServImpl.GetSimulationDeployment(r.Context(), tx, groupID, user)
}

// GetCompetitionRobots returns an array of robots for the competition.
func GetCompetitionRobots(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	return SimServImpl.GetCompetitionRobots(getDefaultApplicationName())
}

// CloudMachineList returns a list with cloud machines (eg. ec2 instances).
// You can request this method with the following cURL request:
//   curl -k -X GET --url http://localhost:8001/1.0/machines
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func CloudMachineList(user *users.User, tx *gorm.DB,
	w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Prepare pagination
	pr, em := ign.NewPaginationRequest(r)
	if em != nil {
		return nil, em
	}

	// Get the parameters
	params := r.URL.Query()
	var status *MachineStatus
	invertStatus := false
	if len(params["status"]) > 0 && len(params["status"][0]) > 0 {
		statusStr := params["status"][0]
		invertStatus = statusStr[0] == '!'
		sliceIndex := 0
		if invertStatus {
			sliceIndex = 1
		}
		status = MachineStatusFrom(statusStr[sliceIndex:])
	}

	var groupID string
	if len(params["groupID"]) > 0 && len(params["groupID"][0]) > 0 {
		groupID = params["groupID"][0]
	}

	// TODO: remove hardcoded application name
	sims, pagination, em := SimServImpl.GetCloudMachineInstances(r.Context(), pr, tx, status, invertStatus, &groupID, user, sptr(getDefaultApplicationName()))
	if em != nil {
		return nil, em
	}

	if pagination != nil {
		ign.WritePaginationHeaders(*pagination, w, r)
	}

	return sims, nil
}

// SimulationWebsocketAddress returns the address of the websocket server and authorization token for the specified
// running simulation. This address and token are meant to be used by the client to establish a websocket connection
// with the simulation to get real-time simulation information.
// A simulation's websocket server is only available while a simulation is running. If a simulation is not running,
// an error will be returned as the websocket server will no longer be available.
// You can request this method with the following curl request:
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/{group}/websocket
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func SimulationWebsocketAddress(user *users.User, tx *gorm.DB, w http.ResponseWriter,
	r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	address, em := SimServImpl.GetSimulationWebsocketAddress(r.Context(), tx, user, groupID)
	if em != nil {
		return nil, em
	}

	return address, nil
}

// SimulationLogGateway returns a URL and a boolean that represent the URL to the proper logs and if it's a file or not.
// If the simulation is running, it will return an URL for live logs and the boolean will be false.
// If the simulation is stopped, it will return an URL for downloadable logs and the boolean will be true.
// You can request this method with the following curl request:
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/{group}/logs/
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func SimulationLogGateway(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	sim, err := GetSimulationDeployment(tx, groupID)

	if err != nil {
		return nil, ign.NewErrorMessage(ign.ErrorUnexpected)
	}

	var logGateway LogGateway
	var path string
	if sim.IsRunning() {
		path = fmt.Sprintf("simulations/%s/logs/live", groupID)
		logGateway = LogGateway{path, false}
	} else {
		path = fmt.Sprintf("simulations/%s/logs/file", groupID)
		logGateway = LogGateway{path, true}
	}

	return logGateway, nil
}

// SimulationLogLive returns a log from a running simulation.
// If the url query includes `lines=N` as parameter,  and the request is for a single simulation
// then this handler will return the last N lines of logs from a live simulation.
// If the url query includes `robot` as parameter with the name of a robot in
// the simulation and the request is for a single simulation, then this will
// return the ROS logs for a specific robot in the simulation.
// If the simulation has been launched as a multisim, it will return a summary of all finished children in the multisim.
// You can request this method with the following curl request:
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/{group}/logs/live?lines=200
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func SimulationLogLive(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	params := r.URL.Query()

	var robotName *string
	if val, ok := params["robot"]; ok {
		robotName = &val[0]
	}

	// TODO: Add environment variable to get default lines value
	lines := int64ptr(1000)
	if val, ok := params["lines"]; ok {
		if n, err := strconv.Atoi(val[0]); err == nil {
			lines = int64ptr(int64(n))
		}
	}

	log, em := SimServImpl.GetSimulationLiveLogs(r.Context(), tx, user, groupID, robotName, lines)

	if em != nil {
		return nil, em
	}

	return log, nil
}

// SimulationLogFileDownload downloads a simulation's logs.
// If the url query includes `link=true` as parameter then this handler will
// return the download URL as a string result instead of doing an http redirect.
// If the url query includes `robot` as parameter with the name of a robot in
// the simulation and the request is for a single simulation, then this will
// return the ROS logs for a specific robot in the simulation.
// You can request this method with the following curl request:
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/{group}/logs/file?link=true
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func SimulationLogFileDownload(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	// Get the parameters
	params := r.URL.Query()
	val, ok := params["link"]
	linkOnly := ok && val[0] == "true"

	var robotName *string
	if val, ok := params["robot"]; ok {
		robotName = &val[0]
	}

	url, em := SimServImpl.GetSimulationLogsForDownload(r.Context(), tx, user, groupID, robotName)
	if em != nil {
		return nil, em
	}

	if linkOnly {
		return *url, nil
	}
	http.Redirect(w, r, *url, http.StatusTemporaryRedirect)
	return nil, nil
}

// QueueGet returns all the simulations from the launch queue
// If the url query includes `page` and `page_size` as parameters then this handler will
// return a paginated list of elements given by those values.
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/queue
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueGet(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	var page *int
	var perPage *int
	params := r.URL.Query()
	if param, ok := params["page"]; ok {
		if value, err := strconv.Atoi(param[0]); err == nil {
			page = intptr(value)
		}
	}
	if param, ok := params["per_page"]; ok {
		if value, err := strconv.Atoi(param[0]); err == nil {
			perPage = intptr(value)
		}
	}

	count, _ := SimServImpl.QueueCount(r.Context(), user)
	w.Header().Set("X-Total-Count", fmt.Sprint(count))

	return SimServImpl.QueueGetElements(r.Context(), user, page, perPage)
}

// QueueCount returns the launch queue elements count
//   curl -k -X GET --url http://localhost:8001/1.0/simulations/queue/count
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueCount(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	return SimServImpl.QueueCount(r.Context(), user)
}

// QueueSwap swaps elements from position A to position B and vice versa
//   curl -k -X PATCH --url http://localhost:8001/1.0/simulations/queue/{groupIDA}/swap/{groupIDB}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueSwap(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupIDA, ok := mux.Vars(r)["groupIDA"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	groupIDB, ok := mux.Vars(r)["groupIDB"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.QueueSwapElements(r.Context(), user, groupIDA, groupIDB)
}

// QueueMoveToFront moves the element to the front of the queue.
//   curl -k -X PATCH --url http://localhost:8001/1.0/simulations/queue/{groupID}/move/front
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueMoveToFront(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.QueueMoveElementToFront(r.Context(), user, groupID)
}

// QueueMoveToBack moves the element to the back of the queue.
//   curl -k -X PATCH --url http://localhost:8001/1.0/simulations/queue/{groupID}/move/back
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueMoveToBack(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.QueueMoveElementToBack(r.Context(), user, groupID)
}

// QueueRemove removes an element from the queue.
//   curl -k -X DELETE --url http://localhost:8001/1.0/simulations/queue/{groupID}
//     --header 'authorization: Bearer <A_VALID_AUTH0_JWT_TOKEN>'
func QueueRemove(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return SimServImpl.QueueRemoveElement(r.Context(), user, groupID)
}

// Healthz returns a string to confirm that cloudsim is running.
//   curl -k -X GET --url http://localhost:8001/1.0/healthz
func Healthz(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	return "Cloudsim is up", nil
}

// Debug is a debug endpoint to get internal state information about a simulation.
func Debug(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	gid := mux.Vars(r)["groupID"]
	return SimServImpl.Debug(user, simulations.GroupID(gid))
}

// ReconnectWebsocket is the endpoint that reconnects simulations to their respective websocket server
func ReconnectWebsocket(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	// CreateSimulation is the input form
	var input ReconnectSimulationList
	if em := ParseStruct(&input, r, true); em != nil {
		return nil, em
	}

	return SimServImpl.ReconnectWebsocket(user, input)
}
