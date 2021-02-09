package simulator

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"testing"
)

func TestStartSimulationAction(t *testing.T) {
	s := NewSimulator(Config{
		DB:                  nil,
		Platform:            nil,
		ApplicationServices: nil,
		ActionService:       nil,
	})

	ctx := context.Background()
	gid := simulations.GroupID(uuid.NewV4().String())

	err := s.Start(ctx, gid)
	assert.NoError(t, err)
}
