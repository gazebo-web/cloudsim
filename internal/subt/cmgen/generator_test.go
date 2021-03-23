package cmgen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestGenerateGazebo(t *testing.T) {
	token := "test-token"
	maxConn := 500
	seed := 5678

	world := "cloudsim_sim.ign;worldName:=tunnel_circuit_practice_01"

	fakeRobotA := fake.NewRobot("testA", "X1")
	fakeRobotB := fake.NewRobot("testB", "X2")

	duration := 1500 * time.Second

	cmd := Gazebo(GazeboConfig{
		World:                   world,
		WorldMaxSimSeconds:      duration,
		Seed:                    &seed,
		AuthorizationToken:      &token,
		MaxWebsocketConnections: maxConn,
		Robots:                  []simulations.Robot{fakeRobotA, fakeRobotB},
		Marsupials: []simulations.Marsupial{
			simulations.NewMarsupial(fakeRobotA, fakeRobotB),
		},
		RosEnabled: true,
	})

	assert.Equal(t, "cloudsim_sim.ign", cmd[0])
	assert.Equal(t, "worldName:=tunnel_circuit_practice_01", cmd[1])
	assert.Equal(t, fmt.Sprintf("durationSec:=%d", int(duration.Seconds())), cmd[2])
	assert.Equal(t, "headless:=true", cmd[3])
	assert.Equal(t, fmt.Sprintf("seed:=%d", seed), cmd[4])
	assert.Equal(t, fmt.Sprintf("websocketAuthKey:=%s", token), cmd[5])
	assert.Equal(t, fmt.Sprintf("websocketAdminAuthKey:=%s", token), cmd[6])
	assert.Equal(t, fmt.Sprintf("websocketMaxConnections:=%d", maxConn), cmd[7])

	assert.Equal(t, fmt.Sprintf("robotName1:=%s", fakeRobotA.GetName()), cmd[8])
	assert.Equal(t, fmt.Sprintf("robotConfig1:=%s", fakeRobotA.GetKind()), cmd[9])

	assert.Equal(t, fmt.Sprintf("robotName2:=%s", fakeRobotB.GetName()), cmd[10])
	assert.Equal(t, fmt.Sprintf("robotConfig2:=%s", fakeRobotB.GetKind()), cmd[11])

	assert.Equal(t, fmt.Sprintf("marsupial1:=%s:%s", fakeRobotA.GetName(), fakeRobotB.GetName()), cmd[12])

	assert.Equal(t, fmt.Sprintf("ros:=%t", true), cmd[13])
}

func TestGenerateCommsBridge(t *testing.T) {
	CommsBridge()
}

func CommsBridge() {

}
