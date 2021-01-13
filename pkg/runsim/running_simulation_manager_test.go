package runsim

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}

type managerTestSuite struct {
	suite.Suite
}
