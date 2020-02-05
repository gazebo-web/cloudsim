package simulations

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimulationDeployment_Clone(t *testing.T) {
	simDep := &SimulationDeployment{
		ID:               123,
		CreatedAt:        time.Now().Add(-time.Hour * 24),
		UpdatedAt:        time.Now().Add(-time.Hour * 24),
		DeletedAt:        timeptr(time.Now().Add(-time.Hour * 12)),
		StoppedAt:        timeptr(time.Now().Add(-time.Hour * 12)),
		ValidFor:         sptr("6h0m0s"),
		Owner:            sptr("Test"),
		Creator:          sptr("TestUser"),
		Private:          boolptr(true),
		StopOnEnd:        boolptr(false),
		Name:             sptr("TestSimDep"),
		Image:            sptr("test"),
		GroupId:          sptr("11111111-1111-1111-1111-111111111111-c-1"),
		ParentGroupId:    sptr("11111111-1111-1111-1111-111111111111-c-1"),
		MultiSim:         2,
		DeploymentStatus: intptr(90),
		ErrorStatus:      sptr("InitializationFailed"),
		Platform:         sptr("subt"),
		Application:      sptr("subt"),
		Extra:            sptr("{}"),
		ExtraSelector:    sptr("Test Circuit"),
		Robots:           sptr("X1,X2"),
	}

	simDepClone := simDep.Clone()
	simDepClone.GroupId = sptr(fmt.Sprintf("%s-r-1", *simDep.GroupId))

	// Check that the model fields have been cleared
	assert.Equal(t, uint(0), simDepClone.ID)
	assert.Nil(t, simDepClone.DeletedAt)
	assert.Nil(t, simDepClone.StoppedAt)

	// Check that the references are copied and can be overwritten
	assert.NotEqual(t, *simDep.GroupId, *simDepClone.GroupId)
}
