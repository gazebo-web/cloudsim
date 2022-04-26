package fake

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	"net/http"
)

// Service is a fake users.Service implementation.
type Service struct {
	*mock.Mock
}

// UserFromJWT returns the User associated to the http request's JWT token.
// This function can return ErrorAuthJWTInvalid if the token cannot be
// read, or ErrorAuthNoUser no user with such identity exists in the DB.
func (f *Service) UserFromJWT(r *http.Request) (*users.User, bool, *ign.ErrMsg) {
	args := f.Called(r)
	return args.Get(0).(*users.User), args.Bool(1), args.Get(2).(*ign.ErrMsg)
}

// VerifyOwner checks if the 'owner' arg is an organization or a user. If the
// 'owner' is an organization, it verifies that the given 'user' arg has the expected
// permission in the organization. If the 'owner' is a user, it verifies that the
// 'user' arg is the same as the owner.
// Dev note: this is an alternative implementation of ign-fuelserver UserService's VerifyOwner.
func (f *Service) VerifyOwner(owner, user string, p per.Action) (bool, *ign.ErrMsg) {
	args := f.Called(owner, user, p)
	return args.Bool(0), args.Get(1).(*ign.ErrMsg)
}

// CanPerformWithRole checks if the 'owner' arg is an organization or a
// user. If the 'owner' is an organization, it verifies that the given 'user' arg
// is authorized to act as the given Role (or above) in the organization.
// If the 'owner' is a user, it verifies that the 'user' arg is the same as
// the owner.
// As a third alternative, if 'owner' is nil then it checks if the 'user' is part
// of the System Admins.
func (f *Service) CanPerformWithRole(owner *string, user string, role per.Role) (bool, *ign.ErrMsg) {
	args := f.Called(owner, user, role)
	return args.Bool(0), args.Get(1).(*ign.ErrMsg)
}

// QueryForResourceVisibility checks the relationship between requestor (user)
// and the resource owner to formulate a database query to determine whether a
// resource is visible to the user
func (f *Service) QueryForResourceVisibility(q *gorm.DB, owner *string, user *users.User) *gorm.DB {
	args := f.Called(q, owner, user)
	return args.Get(0).(*gorm.DB)
}

// IsAuthorizedForResource checks if user has the permission to perform an action on a
// resource.
func (f *Service) IsAuthorizedForResource(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	args := f.Called(user, resource, action)
	return args.Bool(0), args.Get(1).(*ign.ErrMsg)
}

// AddResourcePermission adds a user (or group) permission on a resource
func (f *Service) AddResourcePermission(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	args := f.Called(user, resource, action)
	return args.Bool(0), args.Get(1).(*ign.ErrMsg)
}

// AddScore creates a score entry for a simulation.
func (f *Service) AddScore(groupID *string, competition *string, circuit *string, owner *string, score *float64,
	sources *string) *ign.ErrMsg {
	args := f.Called(groupID, competition, circuit, owner, score, sources)
	return args.Get(0).(*ign.ErrMsg)
}

// IsSystemAdmin returns a bool indicating if the given user is a system admin.
func (f *Service) IsSystemAdmin(user string) bool {
	args := f.Called(user)
	return args.Bool(0)
}

// GetUserFromUsername returns the user database entry from the username
func (f *Service) GetUserFromUsername(username string) (*users.User, *ign.ErrMsg) {
	args := f.Called(username)
	return args.Get(0).(*users.User), args.Get(1).(*ign.ErrMsg)
}

// GetOrganization gets a user's organization database entry from the username
func (f *Service) GetOrganization(username string) (*users.Organization, *ign.ErrMsg) {
	args := f.Called(username)
	return args.Get(0).(*users.Organization), args.Get(1).(*ign.ErrMsg)
}

// StartAutoLoadPolicy starts the auto load remote policy
func (f *Service) StartAutoLoadPolicy() {
	f.Called()
}

// NewFakeService initializes a new fake user service implementation.
// This provider uses the mock library
func NewFakeService() *Service {
	return &Service{
		Mock: new(mock.Mock),
	}
}
