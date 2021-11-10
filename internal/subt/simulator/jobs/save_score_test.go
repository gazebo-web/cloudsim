package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	simsubtfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

// fakeScoreSaver is used to mock the users.Service to submit scores for a specific simulation.
type fakeScoreSaver struct {
	users.Service
}

// AddScore mocks the score saving mechanism on fuel.
func (f *fakeScoreSaver) AddScore(groupID *string, competition *string, circuit *string, owner *string, score *float64, sources *string) *ign.ErrMsg {
	return nil
}

func TestSaveScore(t *testing.T) {
	// Mock simulation service
	simFakeService := simfake.NewService()

	// Create data for test
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd-c-1")
	score := 1.3

	child := simsubtfake.NewSimulation(simsubtfake.SimulationConfig{
		GroupID: gid, Kind: simulations.SimChild, Track: "test",
	})

	parent := simsubtfake.NewSimulation(simsubtfake.SimulationConfig{
		GroupID: "aaaa-bbbb-cccc-dddd", Kind: simulations.SimParent,
	})

	// Mock UpdateScore call
	simFakeService.On("UpdateScore", gid, &score).Return(error(nil))

	// Mock Get call
	simFakeService.On("Get", gid).Return(simulations.Simulation(child), error(nil))

	// Mock GetParent call
	simFakeService.On("GetParent", gid).Return(simulations.Simulation(parent), error(nil))

	// Initialize fake score service
	var fakeScoreService fakeScoreSaver

	// Initialize application services for subt
	appServices := subtapp.NewServices(application.NewServices(simFakeService, &fakeScoreService, nil), nil, nil)

	// Create a new stop state
	startState := state.NewStopSimulation(nil, appServices, gid)

	// Set the score to the state
	startState.Score = &score

	// Initialize a new store with the state
	s := actions.NewStore(startState)

	// Run the job
	_, err := SaveScore.Run(s, nil, nil, startState)

	// Expect no errors when running the save score job
	assert.NoError(t, err)
}
