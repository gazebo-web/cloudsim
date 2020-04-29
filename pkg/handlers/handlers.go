package handlers

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

type handlerWithUser func(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)

// WithUser is a middleware that checks for a valid user from the JWT and passes
// the user to the handlerWithUser.
func WithUser(service users.IService, handler handlerWithUser) ign.HandlerWithResult {
	return func(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
		// Get JWT user. Fail if invalid or missing
		user, ok, em := service.UserFromJWT(r)
		if !ok {
			return nil, em
		}
		return handler(user, w, r)
	}
}

// AfterFn is a middleware that runs a function after executing a given handler.
// It wraps the given handler and returns a HandlerWithResult that will run the given function
// after a successful handler execution.
func AfterFn(h ign.HandlerWithResult, f func(payload interface{}) (interface{}, *ign.ErrMsg)) ign.HandlerWithResult {
	return func(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (i interface{}, msg *ign.ErrMsg) {
		result, err := h(tx, w, r)
		if err != nil {
			return nil, err
		}
		return f(result)
	}
}
