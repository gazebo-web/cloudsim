package simulations

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunningSimulation_IsExpired(t *testing.T) {
	ctx := context.Background()

	dep := SimulationDeployment{
		ID:        0,
		GroupId:   sptr("group-Id-test"),
		Owner:     sptr("Test"),
		CreatedAt: time.Now().Add(-10 * time.Hour),
		ValidFor:  sptr("6h0m0s"),
	}

	rs, err := NewRunningSimulation(ctx, &dep, "test", "test", 720)
	if err != nil {
		t.Fail()
	}
	assert.False(t, rs.IsExpired())
}
