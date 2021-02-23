package nps

// This file defines the controller, which handles route requests. An
// application creates an instance of a controller by calling `NewController`.

import (
	"fmt"
	"github.com/go-playground/form"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// Controller is an interface designed to handle route requests.
type Controller interface {
	Start(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Stop(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	// service is this controller's implementation of the
	// cloudsim/pkg/simulations service. See the simulations_service.go file.
	service     Service
	formDecoder *form.Decoder
}

// NewController creates a new controller
func NewController(db *gorm.DB, logger ign.Logger) Controller {
	return &controller{
		// Create a simulation service to manage simulation instances
		service:     NewService(db, logger),
		formDecoder: form.NewDecoder(),
	}
}

// Builds the ErrMsg extra info from the given DecodeErrors
// \todo: Move this to a common place so that we don't have to copy it
// every time we create an application.
func getDecodeErrorsExtraInfo(err error) []string {
	errs := err.(form.DecodeErrors)
	extra := make([]string, 0, len(errs))
	for field, er := range errs {
		extra = append(extra, fmt.Sprintf("Field: %s. %v", field, er.Error()))
	}
	return extra
}

// Start handles the `/start` route.
//
// Flow: user --> POST /start --> controller.Start()
func (ctrl *controller) Start(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Parse form's values and files.
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	// Get needed data to start simulation from the HTTP request, pass it to the Start Request
	var req StartRequest

	if errs := ctrl.formDecoder.Decode(&req, r.Form); errs != nil {
		fmt.Printf("Failed to decode form")
		return nil, ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs, getDecodeErrorsExtraInfo(errs))
	}

	// Hand off the start request data to the service.
	res, err := ctrl.service.Start(r.Context(), req)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}

	// Remove after addressing next comment
  fmt.Printf("&&&&&&&&&&&&&&&&&&&&&&&&&\n")
	fmt.Println(res)

	// Send response to the user
  return res, nil
}

// Stop handles the `/stop` route.
func (ctrl *controller) Stop(w http.ResponseWriter, r *http.Request) {
	// Parse request

	// Get needed data to stop simulation from the HTTP request, pass it to the Stop Request
	req := StopRequest{}

	res, err := ctrl.service.Stop(r.Context(), req)
	if err != nil {
		// Send error message
	}

	// Remove after addressing next comment
	fmt.Println(res)

	// Send response to the user
}
