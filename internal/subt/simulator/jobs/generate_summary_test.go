package jobs

import (
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

func TestGenerateSummary(t *testing.T) {
	simService := simfake.NewService()
	baseService := application.NewServices(simService, nil)
	s := subtapp.NewServices(baseService, nil, nil)

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := simfake.NewSimulation(gid, simulations.StatusProcessingResults, simulations.SimSingle, nil, "test.org/test")

	stopState := state.NewStopSimulation(nil, s, gid)

	stopState.Stats = simulations.Statistics{
		Started:        1,
		SimulationTime: 2,
		RealTime:       3,
		ModelCount:     4,
	}

	stopState.Score = 5

	store := actions.NewStore(stopState)

	simService.On("Get", gid).Return(sim, error(nil))

	output, err := GenerateSummary.Run(store, nil, nil, stopState)
	require.NoError(t, err)

	resultState, ok := output.(*state.StopSimulation)
	require.True(t, ok)

	assert.NotNil(t, resultState.Summary.GroupID)
	assert.Equal(t, gid, *resultState.Summary.GroupID)

	assert.Equal(t, float64(stopState.Stats.ModelCount), resultState.Summary.ModelCountAvg)
	assert.Equal(t, float64(stopState.Stats.RealTime), resultState.Summary.RealTimeDurationAvg)
	assert.Equal(t, float64(stopState.Stats.SimulationTime), resultState.Summary.SimTimeDurationAvg)

	assert.Equal(t, 0.0, resultState.Summary.SimTimeDurationStdDev)
	assert.Equal(t, 0.0, resultState.Summary.RealTimeDurationStdDev)
	assert.Equal(t, 0.0, resultState.Summary.ModelCountStdDev)

	assert.Equal(t, stopState.Score, resultState.Summary.Score)

}
