package queue

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"strconv"
)

type Controller interface {
	GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Count(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Swap(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	MoveToFront(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	MoveToBack(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Remove(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	var c Controller
	c = &controller{service: service}
	return c
}

func (c *controller) GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	var page *int
	var perPage *int
	params := r.URL.Query()
	if param, ok := params["page"]; ok {
		if value, err := strconv.Atoi(param[0]); err == nil {
			page = tools.Intptr(value)
		}
	}
	if param, ok := params["per_page"]; ok {
		if value, err := strconv.Atoi(param[0]); err == nil {
			perPage = tools.Intptr(value)
		}
	}

	count, _ := c.service.Count(r.Context(), user)
	w.Header().Set("X-Total-Count", fmt.Sprint(count))

	return c.service.GetAll(r.Context(), user, page, perPage)
}

func (c *controller) Count(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	return c.service.Count(r.Context(), user)
}

func (c *controller) Swap(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupIDA, ok := mux.Vars(r)["groupIDA"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	groupIDB, ok := mux.Vars(r)["groupIDB"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	return c.service.Swap(r.Context(), user, groupIDA, groupIDB)
}

func (c *controller) MoveToFront(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.service.MoveToFront(r.Context(), user, groupID)
}

func (c *controller) MoveToBack(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.service.MoveToBack(r.Context(), user, groupID)
}

func (c *controller) Remove(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupID"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.service.Remove(r.Context(), user, groupID)
}
