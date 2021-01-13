package runsim

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRunningSimulation(t *testing.T) {
	suite.Run(t, new(runningSimulationTestSuite))
}

type runningSimulationTestSuite struct {
	suite.Suite
}
