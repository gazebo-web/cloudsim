package users

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// TODO: Implement repository for users.
// Service is the default implementation of IService interface.
type ServiceMock struct {
	UserFromJWTMock func(r *http.Request) (*users.User, bool, *ign.ErrMsg)
	VerifyOwnerMock func(owner, user string, p per.Action) (bool, *ign.ErrMsg)
	CanPerformWithRoleMock func(owner *string, user string, role per.Role) (bool, *ign.ErrMsg)
	QueryForResourceVisibilityMock func(q *gorm.DB, owner *string, user *users.User) *gorm.DB
	IsAuthorizedForResourceMock func(user, resource string, action per.Action) (bool, *ign.ErrMsg)
	AddResourcePermissionMock func(user, resource string, action per.Action) (bool, *ign.ErrMsg)
	AddScoreMock func(groupID *string, competition *string, circuit *string, owner *string, score *float64,
	sources *string) *ign.ErrMsg
	IsSystemAdminMock func(user string) bool
	GetUserFromUsernameMock func(username string) (*users.User, *ign.ErrMsg)
	GetOrganizationMock func(username string) (*users.Organization, *ign.ErrMsg)
}

func NewUserServiceMock() *ServiceMock {
	m := &ServiceMock{}
	return m
}

func (s *ServiceMock) UserFromJWT(r *http.Request) (*users.User, bool, *ign.ErrMsg) {
	return s.UserFromJWTMock(r)
}

func (s *ServiceMock) VerifyOwner(owner, user string, p per.Action) (bool, *ign.ErrMsg) {
	return s.VerifyOwnerMock(owner, user, p)
}

func (s *ServiceMock) CanPerformWithRole(owner *string, user string, role per.Role) (bool, *ign.ErrMsg) {
	return s.CanPerformWithRoleMock(owner, user, role)
}

func (s *ServiceMock) QueryForResourceVisibility(q *gorm.DB, owner *string, user *users.User) *gorm.DB {
	return s.QueryForResourceVisibilityMock(q, owner, user)
}

func (s *ServiceMock) IsAuthorizedForResource(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	return s.IsAuthorizedForResourceMock(user, resource, action)
}

func (s *ServiceMock) AddResourcePermission(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	return s.AddResourcePermissionMock(user, resource, action)
}

func (s *ServiceMock) AddScore(groupID *string, competition *string, circuit *string, owner *string, score *float64,
	sources *string) *ign.ErrMsg {
	return s.AddScoreMock(groupID, competition, circuit, owner, score, sources)
}

func (s *ServiceMock) IsSystemAdmin(user string) bool {
	return s.IsSystemAdminMock(user)
}

func (s *ServiceMock) GetUserFromUsername(username string) (*users.User, *ign.ErrMsg) {
	return s.GetUserFromUsernameMock(username)
}

func (s *ServiceMock) GetOrganization(username string) (*users.Organization, *ign.ErrMsg) {
	return s.GetOrganizationMock(username)
}