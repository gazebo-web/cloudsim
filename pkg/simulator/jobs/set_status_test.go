package jobs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
)

func TestSetStatus(t *testing.T) {
	const gid = simulations.GroupID("aaaa-bbbb-cccc-dddd")
	var expectedStatus = simulations.StatusRunning

	simFakeService := simfake.NewService()

	simFakeService.On("UpdateStatus", gid, expectedStatus).Return(error(nil))

	subtServices := subtapp.NewServices(application.NewServices(simFakeService, nil), nil, nil)

	s := state.NewStartSimulation(nil, subtServices, gid)

	store := actions.NewStore(s)

	result, err := SetSimulationStatus.Run(store, nil, nil, SetSimulationStatusInput{
		GroupID: gid,
		Status:  expectedStatus,
	})
	require.NoError(t, err)

	output, ok := result.(SetSimulationStatusOutput)
	require.True(t, ok)

	assert.Equal(t, expectedStatus, output.Status)
	assert.Equal(t, gid, output.GroupID)
}

func TestSetStatusWithError(t *testing.T) {
	const gid = simulations.GroupID("aaaa-bbbb-cccc-dddd")
	expectedStatus := simulations.StatusRunning
	expectedError := errors.New("test")

	simFakeService := simfake.NewService()

	simFakeService.On("UpdateStatus", gid, expectedStatus).Return(expectedError)

	subtServices := subtapp.NewServices(application.NewServices(simFakeService, nil), nil, nil)

	s := state.NewStartSimulation(nil, subtServices, gid)

	store := actions.NewStore(s)

	_, err := SetSimulationStatus.Run(store, nil, nil, SetSimulationStatusInput{
		GroupID: gid,
		Status:  expectedStatus,
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}
