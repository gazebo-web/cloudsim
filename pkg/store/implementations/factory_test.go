package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/structs"
	"reflect"
	"testing"
)

func TestStoreFactorySuite(t *testing.T) {
	suite.Run(t, new(testStoreFactorySuite))
}

type testStoreFactorySuite struct {
	suite.Suite
}

func (s *testStoreFactorySuite) createFactoryConfig(objectType string,
	config factory.ConfigValues) *factory.Config {

	return &factory.Config{
		Type:   objectType,
		Config: config,
	}
}

func (s *testStoreFactorySuite) TestNewStore() {
	// Prepare config
	key := "test"
	config := factory.Config{
		Type: Store,
		Config: factory.ConfigValues{
			"machinesStore": factory.ConfigValues{
				"keyNameValue": key,
			},
			"ignitionStore": factory.ConfigValues{},
			"orchestratorStore": factory.ConfigValues{
				"ingressNameValue": "test",
			},
		},
	}

	var store store.Store
	s.Equal(nil, Factory.New(&config, nil, &store))
	s.NotEqual(nil, store)
	s.NotEqual(nil, store.Machines())
	s.NotEqual(nil, store.Ignition())
	s.NotEqual(nil, store.Orchestrator())

	// Validate the type of the returned object
	s.Equal("*store.store", reflect.TypeOf(store).String())

	// Validate provided value
	s.Equal(key, store.Machines().KeyName())

	// Validate default value
	machineType, err := structs.GetFieldTagValue(store.Machines(), "MachineTypeValue", "default")
	s.Equal(nil, err)
	s.Equal(machineType, store.Machines().Type())
}
