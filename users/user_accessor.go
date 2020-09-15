package users

import (
	"context"
	"github.com/caarlos0/env"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/subt"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"strings"
	"time"
)

type userAccessorConf struct {
	AutoLoadPolicySeconds int `env:"USER_ACCESSOR_AUTOLOAD_SECONDS" envDefault:"10"`
	sysAdmin              string
}

// UserAccessor is used by the cloudsim server to remotely get Users and their membership
// to Organizations.
type UserAccessor interface {
	// UserFromJWT returns the User associated to the http request's JWT token.
	// This function can return ErrorAuthJWTInvalid if the token cannot be
	// read, or ErrorAuthNoUser no user with such identity exists in the DB.
	UserFromJWT(r *http.Request) (*users.User, bool, *ign.ErrMsg)
	// VerifyOwner checks if the 'owner' arg is an organization or a user. If the
	// 'owner' is an organization, it verifies that the given 'user' arg has the expected
	// permission in the organization. If the 'owner' is a user, it verifies that the
	// 'user' arg is the same as the owner.
	// Dev note: this is an alternative implementation of ign-fuelserver UserService's VerifyOwner.
	VerifyOwner(owner, user string, p per.Action) (bool, *ign.ErrMsg)
	// CanPerformWithRole checks if the 'owner' arg is an organization or a
	// user. If the 'owner' is an organization, it verifies that the given 'user' arg
	// is authorized to act as the given Role (or above) in the organization.
	// If the 'owner' is a user, it verifies that the 'user' arg is the same as
	// the owner.
	// As a third alternative, if 'owner' is nil then it checks if the 'user' is part
	// of the System Admins.
	CanPerformWithRole(owner *string, user string, role per.Role) (bool, *ign.ErrMsg)
	// QueryForResourceVisibility checks the relationship between requestor (user)
	// and the resource owner to formulate a database query to determine whether a
	// resource is visible to the user
	QueryForResourceVisibility(q *gorm.DB, owner *string, user *users.User) *gorm.DB
	// IsAuthorizedForResource checks if user has the permission to perform an action on a
	// resource.
	IsAuthorizedForResource(user, resource string, action per.Action) (bool, *ign.ErrMsg)
	// AddResourcePermission adds a user (or group) permission on a resource
	AddResourcePermission(user, resource string, action per.Action) (bool, *ign.ErrMsg)
	// AddScore creates a score entry for a simulation.
	AddScore(groupID *string, competition *string, circuit *string, owner *string, score *float64,
		sources *string) *ign.ErrMsg
	// IsParticipant returns true if the owner is a participant in the competition.
	IsParticipant(competition, owner string) bool
	// IsSystemAdmin returns a bool indicating if the given user is a system admin.
	IsSystemAdmin(user string) bool
	// GetUserFromUsername returns the user database entry from the username
	GetUserFromUsername(username string) (*users.User, *ign.ErrMsg)
	// GetOrganization gets a user's organization database entry from the username
	GetOrganization(username string) (*users.Organization, *ign.ErrMsg)
}

// UserAccessorImpl is the default implementation of UserAccessor interface.
type UserAccessorImpl struct {
	// The UserAccessor config. Read from environment variables
	cfg userAccessorConf
	// Global database interface to Users DB
	Db *gorm.DB
	// Membership and permissions for Users/Orgs.
	p              *per.Permissions
	syncedEnforcer *casbin.SyncedEnforcer
	// access to permissions over resources (not users/orgs) membership.
	resourcePermissions *per.Permissions
}

// NewUserAccessor initializes a new UserAccessor.
func NewUserAccessor(ctx context.Context, resourcePermissions *per.Permissions, usersDb *gorm.DB, sysAdmin string) (*UserAccessorImpl, error) {

	ua := UserAccessorImpl{}
	ua.Db = usersDb
	ua.resourcePermissions = resourcePermissions

	// Read configuration from environment
	ua.cfg = userAccessorConf{}
	if err := env.Parse(&ua.cfg); err != nil {
		return nil, err
	}
	ua.cfg.sysAdmin = sysAdmin

	// Create Casbin helpers
	adapter := gormadapter.NewAdapterByDB(usersDb)
	enforcer := casbin.NewSyncedEnforcer("permissions/policy.conf", adapter)
	ua.syncedEnforcer = enforcer

	ua.p = &per.Permissions{}
	ua.p.InitWithEnforcerAndAdapter(enforcer, adapter, sysAdmin)

	return &ua, nil
}

// StartAutoLoadPolicy starts the auto load remote policy
func (u *UserAccessorImpl) StartAutoLoadPolicy() {
	// Auto load remote policy
	u.syncedEnforcer.StartAutoLoadPolicy(time.Duration(u.cfg.AutoLoadPolicySeconds) * time.Second)
}

// UserFromJWT returns the User associated to the http request's JWT token.
// This function can return ErrorAuthJWTInvalid if the token cannot be
// read, or ErrorAuthNoUser no user with such identity exists in the DB.
func (u *UserAccessorImpl) UserFromJWT(r *http.Request) (*users.User, bool, *ign.ErrMsg) {
	return getUserFromToken(u.Db, r)
}

// getUserFromToken returns the User associated to the http request's JWT token.
// This function can return ErrorAuthJWTInvalid if the token cannot be
// read, or ErrorAuthNoUser no user with such identity exists in the DB.
func getUserFromToken(tx *gorm.DB, r *http.Request) (*users.User, bool, *ign.ErrMsg) {
	var user *users.User
	if token := r.Header.Get("Private-Token"); len(token) > 0 {
		var accessToken *ign.AccessToken
		var err *ign.ErrMsg
		if accessToken, err = ign.ValidateAccessToken(token, tx); err != nil {
			return nil, false, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}

		user = new(users.User)
		if err := tx.Where("id = ?", accessToken.UserID).First(user).Error; err != nil {
			return nil, false, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
	} else {
		identity, valid := ign.GetUserIdentity(r)
		if !valid {
			return nil, false, ign.NewErrorMessage(ign.ErrorAuthJWTInvalid)
		}

		var em *ign.ErrMsg
		user, em = users.ByIdentity(tx, identity, false)
		if em != nil {
			return nil, false, em
		}
	}

	return user, true, nil
}

// VerifyOwner checks to see if the 'owner' arg is an organization or a user. If the
// 'owner' is an organization, it verifies that the given 'user' arg has the expected
// permission in the organization. If the 'owner' is a user, it verifies that the
// 'user' arg is the same as the owner.
// Dev note: this is an alternative implementation of ign-fuelserver UserService's VerifyOwner.
func (u *UserAccessorImpl) VerifyOwner(owner, user string, p per.Action) (bool, *ign.ErrMsg) {
	// check if owner is an organization
	org, em := users.ByOrganizationName(u.Db, owner, false)
	if org != nil && em == nil {
		// check if user has at least the given permission in that organization
		ok, em := u.p.IsAuthorized(user, *org.Name, p)
		if !ok {
			return false, em
		}
	} else {
		// Owner is a user. Make sure the owner is the same as the jwt user.
		if owner != user {
			// jwt user is different from owner field!
			return false, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
	}
	return true, nil
}

// CanPerformWithRole checks to see if the 'owner' arg is an organization or a
// user. If the 'owner' is an organization, it verifies that the given 'user' arg
// is authorized to act as the given Role (or above) in the organization.
// If the 'owner' is a user, it verifies that the 'user' arg is the same as
// the owner.
// As a third alternative, if 'owner' is nil then it checks if the 'user' is part
// of the System Admins.
func (u *UserAccessorImpl) CanPerformWithRole(owner *string, user string, role per.Role) (bool, *ign.ErrMsg) {
	if owner == nil {
		ok := u.p.IsSystemAdmin(user)
		if !ok {
			return false, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
		return true, nil
	}

	// check if owner is an organization
	org, em := users.ByOrganizationName(u.Db, *owner, false)
	if org != nil && em == nil {
		// check if user can act with the given role in the organization
		ok, em := u.p.IsAuthorizedForRole(user, *org.Name, role)
		if !ok {
			return false, em
		}
	} else {
		// Owner is a user. Make sure the owner is the same as the jwt user.
		if *owner != user {
			return false, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
	}
	return true, nil
}

// QueryForResourceVisibility checks the relationship between requestor (user)
// and the resource owner to formulate a database query to determine whether a
// resource is visible to the user
func (u *UserAccessorImpl) QueryForResourceVisibility(q *gorm.DB, owner *string, user *users.User) *gorm.DB {
	// Check resource visibility
	publicOnly := false
	// if owner is specified
	if owner != nil {
		if user == nil {
			// if no user is specified, only public resources are visible
			publicOnly = true
		} else {
			// check if owner is an org
			org, _ := users.ByOrganizationName(u.Db, *owner, false)
			if org != nil {
				// if owner is an org, check if requestor is part of that org
				ok, _ := u.p.IsAuthorized(*user.Username, *org.Name, per.Read)
				if !ok {
					// if requestor is not part of that org, only public resources will
					// be returned
					publicOnly = true
				}
			} else if *user.Username != *owner {
				// if owner is not an org then this is another user's resource
				// TODO check permissions when resource sharing is implemented
				// but for now assume user can only acccess other user's public
				// resources
				publicOnly = true
			}
		}
		if !publicOnly {
			q = q.Where("owner = ?", *owner)
		} else {
			q = q.Where("owner = ? AND private = ?", *owner, 0)
		}
	} else {
		// if owner is not specified, the query should only return resources that
		// are either 1) public or 2) private resources that requestor has read
		// permissions
		if user == nil {
			q = q.Where("private = ?", 0)
		} else {
			userGroups := u.p.GetGroupsForUser(*user.Username)
			userGroups = append(userGroups, *user.Username)
			q = q.Where("private = ? OR owner IN (?)", 0, userGroups)
		}
	}
	return q
}

// IsAuthorizedForResource checks if user has the permission to perform an action on a
// resource.
func (u *UserAccessorImpl) IsAuthorizedForResource(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	ok, _ := u.resourcePermissions.IsAuthorized(user, resource, action)
	if ok {
		return true, nil
	}

	// Get the groups to which the user belongs and check again
	userGroups := u.p.GetGroupsForUser(user)
	for _, g := range userGroups {
		ok, _ := u.resourcePermissions.IsAuthorized(g, resource, action)
		if ok {
			return true, nil
		}
	}

	return false, ign.NewErrorMessage(ign.ErrorUnauthorized)
}

// AddResourcePermission adds a user (or group) permission on a resource
func (u *UserAccessorImpl) AddResourcePermission(user, resource string, action per.Action) (bool, *ign.ErrMsg) {
	return u.resourcePermissions.AddPermission(user, resource, action)
}

// IsParticipant returns true if the owner is a participant in the competition.
func (u *UserAccessorImpl) IsParticipant(competition, owner string) bool {
	var participant subt.CompetitionParticipant
	if err := u.Db.Where("competition = ? AND owner = ?", competition, owner).First(&participant).Error; err != nil {
		return false
	}

	return true
}

// AddScore creates a new score entry for an owner in a competition circuit
// TODO HACK This is accessing Fuel's database directly
func (u *UserAccessorImpl) AddScore(groupID *string, competition *string, circuit *string, owner *string,
	score *float64, sources *string) *ign.ErrMsg {
	entry := subt.CompetitionScore{
		GroupID:     groupID,
		Competition: competition,
		Circuit:     circuit,
		Owner:       owner,
		Score:       score,
		Sources:     sources,
	}
	if err := u.Db.Create(&entry).Error; err != nil {
		ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	return nil
}

// IsSystemAdmin returns a bool indicating if the given user is a system admin.
func (u *UserAccessorImpl) IsSystemAdmin(user string) bool {
	return u.resourcePermissions.IsSystemAdmin(user)
}

// GetUserFromUsername gets the user database entry from the username
func (u *UserAccessorImpl) GetUserFromUsername(username string) (*users.User, *ign.ErrMsg) {
	user := &users.User{}
	if err := u.Db.
		Model(user).
		Where("LOWER(username) = ?", strings.ToLower(username)).
		First(user).
		Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

	return user, nil
}

// GetOrganization gets a user's organization database entry from the username
func (u *UserAccessorImpl) GetOrganization(name string) (*users.Organization, *ign.ErrMsg) {
	org := &users.Organization{}
	if err := u.Db.
		Model(org).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		First(org).
		Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

	return org, nil
}
