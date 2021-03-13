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
			"ignitionStore": factory.ConfigValues{
				"defaultSenderValue": "test@ignitionrobotics.org",
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
