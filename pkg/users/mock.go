package users

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// TODO: Implement repository for users.
// service is the default implementation of Service interface.
type ServiceMock struct {
	*mock.Mock
}

func NewServiceMock() *ServiceMock {
	m := &ServiceMock{
		Mock: new(mock.Mock),
	}
	return m
}

func (s *ServiceMock) UserFromJWT(r *http.Request) (*users.User, bool, *ign.ErrMsg) {
	args := s.Called(r)

	var user *users.User
	user = args.Get(0).(*users.User)

	var err *ign.ErrMsg
	err = args.Get(2).(*ign.ErrMsg)

	return user, args.Bool(1), err
}

func (s *ServiceMock) VerifyOwner(owner, user string, p per.Action) (bool, *ign.ErrMsg) {
	args := s.Called(owner, user, p)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Bool(0), err
}

func (s *ServiceMock) CanPerformWithRole(owner *string, user string, role per.Role) (bool, *ign.ErrMsg) {
	args := s.Called(owner, user, role)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Bool(0), err
}

func (s *ServiceMock) QueryForResourceVisibility(q *gorm.DB, owner *string, user *users.User) *gorm.DB {
	args := s.Called(q, owner, user)

	var db *gorm.DB
	db = args.Get(0).(*gorm.DB)

	return db
}

func (s *ServiceMock) IsAuthorizedForResource(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	args := s.Called(user, resource, action)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Bool(0), err
}

func (s *ServiceMock) AddResourcePermission(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	args := s.Called(user, resource, action)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Bool(0), err
}

func (s *ServiceMock) AddScore(groupID *string, competition *string, circuit *string, owner *string, score *float64,
	sources *string) *ign.ErrMsg {
	args := s.Called(groupID, competition, circuit, owner, score)

	var err *ign.ErrMsg
	err = args.Get(0).(*ign.ErrMsg)

	return err
}

func (s *ServiceMock) IsSystemAdmin(user string) bool {
	args := s.Called(user)

	return args.Bool(0)
}

func (s *ServiceMock) GetUserFromUsername(username string) (*users.User, *ign.ErrMsg) {
	args := s.Called(username)

	var user *users.User
	user = args.Get(0).(*users.User)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return user, err
}

func (s *ServiceMock) GetOrganization(username string) (*users.Organization, *ign.ErrMsg) {
	args := s.Called(username)

	var org *users.Organization
	org = args.Get(0).(*users.Organization)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return org, err
}
