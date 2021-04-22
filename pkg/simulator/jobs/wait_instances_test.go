package jobs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"testing"
)

func (w *machinesWaiter) WaitOK(input []cloud.WaitMachinesOKInput) error {
	for _, in := range input {
		w.WaitedFor = append(w.WaitedFor, in.Instances...)
	}
	return nil
}

type machinesWaiter struct {
	WaitedFor []string
	cloud.Machines
}

func TestWaitForInstances(t *testing.T) {
	input := WaitForInstancesInput([]cloud.CreateMachinesOutput{
		{
			Instances: []string{
				"i-1234",
				"i-1234",
				"i-1234",
				"i-1234",
				"i-1234",
			},
		},
		{
			Instances: []string{
				"i-1234",
				"i-1234",
			},
		},
	})

	m := machinesWaiter{}

	p := platform.NewPlatform(platform.Components{Machines: &m})
	s := state.NewStartSimulation(p, nil, "aaaa-bbbb-cccc-dddd")
	store := actions.NewStore(s)

	_, err := WaitForInstances.Run(store, nil, nil, input)
	require.NoError(t, err)
	assert.Len(t, m.WaitedFor, 7)
}
