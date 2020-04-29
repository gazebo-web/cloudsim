package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// New initializes a new mux router.
func New() *mux.Router {
	r := ign.NewRouter()
	return r
}

// ConfigureRoutes attaches a set of routes in the given router to run on the server.
// It appends the version and the namespace to those routes.
// Returns the configured router.
func ConfigureRoutes(server *ign.Server, version string, namespace string, routes ign.Routes) {
	prefix := fmt.Sprintf("/%s/%s", version, namespace)
	sub := server.Router.PathPrefix(prefix).Subrouter()
	server.ConfigureRouterWithRoutes(prefix, sub, routes)
}
