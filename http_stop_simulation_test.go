package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/migrations"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPStopSimulationSuite(t *testing.T) {
	suite.Run(t, new(HTTPStopSimulationSuite))
}

type HTTPStopSimulationSuite struct {
	suite.Suite
	DB                *gorm.DB
	ResponseRecorder  *httptest.ResponseRecorder
	UserTeamA         users.User
	UserTeamB         users.User
	UserAdmin         users.User
	UserService       useracc.Service
	Permissions       *permissions.Permissions
	SimulationService *sim.Service
	Router            *mux.Router
}

func (s *HTTPStopSimulationSuite) SetupSuite() {
	db, err := gormUtils.GetTestDBFromEnvVars()
	s.Require().NoError(err)
	s.DB = db

	s.UserTeamA = users.User{
		Username: sptr("user-a"),
	}
	s.UserTeamB = users.User{
		Username: sptr("user-b"),
	}
	s.UserAdmin = users.User{
		Username: sptr("sysadmin"),
	}
}

func (s *HTTPStopSimulationSuite) TearDownSuite() {
	db := s.DB.DB()
	s.Require().NoError(db.Close())
}

func (s *HTTPStopSimulationSuite) SetupTest() {
	ctx := context.Background()
	migrations.DBDropModels(ctx, s.DB)
	migrations.DBMigrate(ctx, s.DB)

	var err error

	s.Permissions = &permissions.Permissions{}
	s.Require().NoError(s.Permissions.Init(s.DB, *s.UserAdmin.Username))

	s.UserService, err = useracc.NewService(ctx, s.Permissions, s.DB, *s.UserAdmin.Username)
	s.Require().NoError(err)

	sim.SimServImpl, err = sim.NewSimulationsService(
		context.Background(),
		s.DB,
		sim.SynchronicPoolFactory,
		s.UserService,
		true,
	)

	var ok bool
	s.SimulationService, ok = sim.SimServImpl.(*sim.Service)
	s.Require().True(ok)

	s.Require().NoError(err)

	s.Router = mux.NewRouter()

}

func (s *HTTPStopSimulationSuite) TestUserRequestsTerminationWhenEmptyGroupID() {
	req := httptest.NewRequest(
		http.MethodDelete,
		"https://cloudsim.ignitionrobotics.com/1.0/simulations/",
		nil,
	)

	res := httptest.NewRecorder()

	s.Router.HandleFunc("/simulations/{group}", func(w http.ResponseWriter, r *http.Request) {
		_, em := sim.CloudsimSimulationDelete(&s.UserTeamA, s.DB, w, r)
		s.Assert().NotNil(em)
		s.Assert().Error(em.BaseError)

		s.Assert().Equal(ign.ErrorIDNotInRequest, em.ErrCode)
	})

	s.Router.ServeHTTP(res, req)
}

func (s *HTTPStopSimulationSuite) TestUserRequestsTerminationWhenSimulationDoesNotExist() {
	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("https://cloudsim.ignitionrobotics.com/1.0/simulations/%s", "aaaa-bbbb-cccc-dddd"),
		nil,
	)

	res := httptest.NewRecorder()

	s.Router.HandleFunc("/simulations/{group}", func(w http.ResponseWriter, r *http.Request) {
		_, em := sim.CloudsimSimulationDelete(&s.UserTeamA, s.DB, nil, req)

		s.Assert().NotNil(em)
		s.Assert().Error(em.BaseError)
		s.Assert().Equal(ign.ErrorSimGroupNotFound, em.ErrCode)
	})

	s.Router.ServeHTTP(res, req)
}

func (s *HTTPStopSimulationSuite) TestUserTeamAShouldNotStopSimulationFromTeamB() {
	//simulation, err := s.SimulationService.ServiceAdaptor.Create(simulations.CreateSimulationInput{
	//	Name:      "",
	//	Owner:     nil,
	//	Creator:   s.UserAdmin,
	//	Image:     nil,
	//	Private:   false,
	//	StopOnEnd: false,
	//	Extra:     "",
	//	Track:     "",
	//	Robots:    "",
	//})
	//s.Require().NoError(err)
}

func (s *HTTPStopSimulationSuite) TestUserCannotStopCircuitSimulation() {

}

func (s *HTTPStopSimulationSuite) TestAdminAllowedToStopSimulation() {

}

func (s *HTTPStopSimulationSuite) TestAdminAllowedToStopCircuitSimulation() {

}
