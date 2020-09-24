package simulations

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/cmd/subt/transport/ign"
	"testing"
	"time"
)

func TestRunningSimulation_IsExpired(t *testing.T) {
	ctx := context.Background()

	dep := SimulationDeployment{
		ID:        0,
		GroupID:   sptr("group-Id-test"),
		Owner:     sptr("Test"),
		CreatedAt: time.Now().Add(-10 * time.Hour),
		ValidFor:  sptr("6h0m0s"),
	}

	transport := ignws.NewPubSubTransporterMock()

	cb := mock.AnythingOfType("ign.Callback")

	transport.On("Subscribe", "test", cb).Twice().Return(nil)

	rs, err := NewRunningSimulation(ctx, &dep, transport, "test", "test", 720)
	if err != nil {
		t.Fail()
	}
	assert.False(t, rs.IsExpired())
}
