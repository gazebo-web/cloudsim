package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Config struct {
	Version string
}

func New() *mux.Router {
	r := ign.NewRouter()
	return r
}

func ConfigureRoutes(server *ign.Server, router *mux.Router, version string, namespace string, routes ign.Routes) *mux.Router {
	prefix := fmt.Sprintf("/%s/%s", version, namespace)
	sub := router.PathPrefix(prefix).Subrouter()
	server.ConfigureRouterWithRoutes(prefix, sub, routes)
	return router
}