package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	utilsgorm "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"testing"
)

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

type RepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	Repository Repository
}

type Test struct {
	gorm.Model
	Name string `json:"name"`
}

func (suite *RepositoryTestSuite) SetupSuite() {
	db, err := utilsgorm.GetTestDBFromEnvVars()
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.Require().NoError(suite.db.DropTableIfExists(&Test{}).Error)
	suite.Require().NoError(suite.db.AutoMigrate(&Test{}).Error)
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	suite.Require().NoError(suite.db.Close())
}

func (suite *RepositoryTestSuite) TestImplementsInterface() {
	var expected *Repository
	suite.Assert().Implements(expected, new(repositorySQL))
}
