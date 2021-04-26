package jobs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"testing"
)

func TestLaunchInstances(t *testing.T) {
	const (
		gid = simulations.GroupID("aaaa-bbbb-cccc-dddd")
	)

	db, err := gorm.GetTestDBFromEnvVars()
	require.NoError(t, err)

	var (
		m                 = fake.NewMachines()
		p                 = platform.NewPlatform(platform.Components{Machines: m})
		s                 = state.NewStartSimulation(p, nil, gid)
		store             = actions.NewStore(s)
		expectedInstances = []string{"instance-test-a", "instance-test-b"}
		dep               = actions.Deployment{}
	)

	input := []cloud.CreateMachinesInput{
		{},
	}

	m.On("Create", input).Return([]cloud.CreateMachinesOutput{{
		Instances: expectedInstances,
	}}, error(nil))

	result, err := LaunchInstances.Run(store, db, &dep, LaunchInstancesInput(input))
	require.NoError(t, err)

	output, ok := result.(LaunchInstancesOutput)
	require.True(t, ok)

	assert.Len(t, output, 1)
	assert.Equal(t, output[0].Instances, expectedInstances)
}
