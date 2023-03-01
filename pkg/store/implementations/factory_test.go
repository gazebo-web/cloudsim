package implementations

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/store"
	"github.com/gazebo-web/gz-go/v7/structs"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

func TestStoreFactorySuite(t *testing.T) {
	suite.Run(t, new(testStoreFactorySuite))
}

type testStoreFactorySuite struct {
	suite.Suite
}

func (s *testStoreFactorySuite) TestNewStore() {
	// Prepare config
	key := "test"
	config := factory.Config{
		Type: Store,
		Config: factory.ConfigValues{
			"machinesStore": factory.ConfigValues{
				"instanceProfileValue": "test",
				"keyNameValue":         "test",
				"namePrefixValue":      "test",
				"clusterNameValue":     "test",
			},
			"ignitionStore": factory.ConfigValues{
				"defaultSenderValue": "test@ignitionrobotics.org",
				"logsBucketValue":    "test_bucket",
			},
			"orchestratorStore": factory.ConfigValues{
				"ingressNameValue": "test",
				"ingressHostValue": "test.com",
			},
		},
	}

	var store store.Store
	err := Factory.New(&config, nil, &store)
	if err != nil {
		s.FailNow(err.Error())
	}
	s.Require().NotNil(store)
	s.Require().NotNil(store.Machines())
	s.Require().NotNil(store.Ignition())
	s.Require().NotNil(store.Orchestrator())

	// Validate the type of the returned object
	s.Equal("*store.store", reflect.TypeOf(store).String())

	// Validate provided value
	s.Equal(key, store.Machines().KeyName())

	// Validate default value
	machineType, err := structs.GetFieldTagValue(store.Machines(), "MachineTypeValue", "default")
	s.Require().Nil(err)
	s.Equal(machineType, store.Machines().Type())
}
