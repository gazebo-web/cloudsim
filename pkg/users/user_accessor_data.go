package users

import (
	"context"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuelusers "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
	"strings"
)

// Test data and db functions to use with tests, to mock the real Users db.

/*
	It will create a set of default users and organizations used
	during testing. The users and orgs are:
	- SystemAdmin user: the name passed as argument.
	- Org: <Application>. It will create an Org with the passed in application name.
	---- Users:
	------ "AppOrgAdmin", an admin of the <Application> org.
	------ "AppOrgMember", a member of the <Application> org. It will still be an
					"admin" of the competition.
	- Org: "TeamA".
	---- Users:
	------ "TeamAOwner" (ie. username), with Owner role in the org.
	------ "TeamAAdmin", with admin role in the org.
	------ "TeamAUser1", member of the org.
	------ "TeamAUser2", member of the org.
	- Org: "TeamB"
	------ "TeamBOwner", with owner role in the org.
	------ "TeamBAdmin", with admin role in the org.
	------ "TeamBUser1", member of the org.
	------ "TeamBUser2", member of the org.

	Owners and Admins can Read and Write in the Org. Members can only Read.

	Note: the identities of these Users are their usernames.
*/

// TODO these methods should be part of ign-fuelserver to avoid duplicating this initialization code.

// UserAccessorDataMock allows us to configure the IUserAccessor with mock data used in tests.
type UserAccessorDataMock struct {
	ua              *service
	sysadminIdentiy string
	app             string
}

// NewUserAccessorDataMock ...
func NewUserAccessorDataMock(ctx context.Context, ua Service, sysadminIdentiy, application string) *UserAccessorDataMock {
	useracc := ua.(*service)
	mock := UserAccessorDataMock{
		ua:              useracc,
		sysadminIdentiy: sysadminIdentiy,
		app:             application,
	}
	return &mock
}

// ReloadEverything ...
func (m *UserAccessorDataMock) ReloadEverything(ctx context.Context) *ign.ErrMsg {
	// Foreign key checks are temporarily disabled to be able to drop tables
	m.ua.Db.Exec("SET FOREIGN_KEY_CHECKS=0;")

	m.usersDBDropModels(ctx)
	m.usersDBMigrate(ctx)
	if em := m.addTestData(ctx); em != nil {
		return em
	}

	// Foreign key checks are reenabled
	m.ua.Db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	return nil
}

func (m *UserAccessorDataMock) addTestData(ctx context.Context) *ign.ErrMsg {

	usersDb := m.ua.Db

	sysAdminUser := m.createUser(m.ua.cfg.sysAdmin)
	sysAdminUser.Identity = &m.sysadminIdentiy
	if em := m.addUserToDb(usersDb, sysAdminUser); em != nil {
		return em
	}

	appOrg := &fuelusers.Organization{Name: &m.app, Description: &m.app,
		Email: tools.Sptr(m.app + "@email.com"), Creator: sysAdminUser.Username}
	appAdmin := m.createUser(m.app + "Admin")
	if em := m.addUserToDb(usersDb, appAdmin); em != nil {
		return em
	}
	ok, em := m.ua.p.AddUserGroupRole(*appAdmin.Username, *appOrg.Name, per.Admin)
	if !ok {
		return em
	}
	appMember := m.createUser(m.app + "Member")
	if em := m.addUserToDb(usersDb, appMember); em != nil {
		return em
	}
	ok, em = m.ua.p.AddUserGroupRole(*appMember.Username, *appOrg.Name, per.Member)
	if !ok {
		return em
	}
	if em = m.addOrgToDB(usersDb, appOrg, sysAdminUser); em != nil {
		return em
	}

	teams := []string{"TeamA", "TeamB"}
	for _, t := range teams {
		org, users := m.createOrgAndUsers(t)
		for _, u := range users {
			if em := m.addUserToDb(usersDb, u); em != nil {
				return em
			}
			// Add the user to the org. We do this by updating the permissions.
			role := m.getUserRole(*u.Username)
			ok, em := m.ua.p.AddUserGroupRole(*u.Username, *org.Name, *role)
			if !ok {
				return em
			}
		}
		if em := m.addOrgToDB(usersDb, org, users[0]); em != nil {
			return em
		}
	}

	return nil
}

func (m *UserAccessorDataMock) addUserToDb(tx *gorm.DB, u *fuelusers.User) *ign.ErrMsg {
	// Add the user tso the dat-abase.
	// Note: we also need to add (before) a row to UniqueOwners
	owner := fuelusers.UniqueOwner{Name: u.Username, OwnerType: fuelusers.OwnerTypeUser}
	if err := tx.Create(&owner).Create(&u).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	return nil
}

func (m *UserAccessorDataMock) addOrgToDB(tx *gorm.DB, org *fuelusers.Organization, creator *fuelusers.User) *ign.ErrMsg {
	// Create the organization in the permissions db as a 'group' and set the
	// creator as the 'owner'.
	// This is the same as adding the user to the 'default' team of the Org.
	ok, em := m.ua.p.AddUserGroupRole(*creator.Username, *org.Name, per.Owner)
	if !ok {
		log.Println("Error adding owner permission when creating Org", em)
	}

	owner := fuelusers.UniqueOwner{Name: org.Name, OwnerType: fuelusers.OwnerTypeOrg}
	if err := tx.Create(&owner).Create(org).Error; err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	return nil
}

func (m *UserAccessorDataMock) getUserRole(username string) *per.Role {

	if m.ua.cfg.sysAdmin == username {
		role := per.SystemAdmin
		return &role
	}

	if strings.HasSuffix(username, "Owner") {
		role := per.Owner
		return &role
	}
	if strings.HasSuffix(username, "Admin") {
		role := per.Admin
		return &role
	}
	role := per.Member
	return &role
}

func (m *UserAccessorDataMock) createOrgAndUsers(teamName string) (*fuelusers.Organization, []*fuelusers.User) {
	// create the org admin user
	ownerName := teamName + "Owner"
	owner := m.createUser(ownerName)
	// Create the Org mock
	org := &fuelusers.Organization{Name: &teamName, Description: &teamName,
		Email: tools.Sptr(teamName + "@email.com"), Creator: &ownerName}
	// create extra users
	users := []*fuelusers.User{
		owner,
		m.createUser(teamName + "Admin"),
		m.createUser(teamName + "User1"),
		m.createUser(teamName + "User2"),
	}
	return org, users
}

func (m *UserAccessorDataMock) createUser(username string) *fuelusers.User {
	user := &fuelusers.User{
		Identity: &username,
		Name:     &username,
		Username: &username,
		Email:    tools.Sptr("test@email.com"),
	}
	return user
}

// usersDBMigrate creates the users db tables
func (m *UserAccessorDataMock) usersDBMigrate(ctx context.Context) {

	usersDb := m.ua.Db
	p := m.ua.p
	sysAdmin := m.ua.cfg.sysAdmin

	usersDb.AutoMigrate(
		&fuelusers.UniqueOwner{},
		&fuelusers.User{}, &fuelusers.Organization{}, &fuelusers.Team{},
		p.DBTable(),
	)

	// After removing tables we can ask casbin to re initialize
	if err := p.Reload(sysAdmin); err != nil {
		log.Fatal("Error reloading casbin policies", err)
	}

}

// usersDBDropModels drops all tables from db. Used by tests.
func (m *UserAccessorDataMock) usersDBDropModels(ctx context.Context) {

	usersDb := m.ua.Db
	p := m.ua.p

	usersDb.DropTableIfExists(
		&fuelusers.Team{},
		&fuelusers.Organization{},
		&fuelusers.User{},
		&fuelusers.UniqueOwner{},
		p.DBTable(),
	)
}
