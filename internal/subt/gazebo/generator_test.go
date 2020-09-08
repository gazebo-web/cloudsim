package gazebo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	worldIndex := 1
	runIndex := 1
	token := "test-token"
	maxConn := 500
	seeds := []int{1234, 5678, 91011}

	worlds := []string{
		"cloudsim_sim.ign;worldName:=tunnel_circuit_practice_01,",
		"cloudsim_sim.ign;worldName:=tunnel_circuit_practice_02",
		"cloudsim_sim.ign;worldName:=tunnel_circuit_practice_03",
	}

	fakeRobotA := fake.NewRobot("testA", "X1")
	fakeRobotB := fake.NewRobot("testB", "X2")

	duration := 1500 * time.Second

	cmd := Generate(LaunchConfig{
		Worlds:                  worlds,
		WorldMaxSimSeconds:      duration,
		Seeds:                   seeds,
		WorldIndex:              &worldIndex,
		RunIndex:                &runIndex,
		AuthorizationToken:      &token,
		MaxWebsocketConnections: maxConn,
		Robots:                  []simulations.Robot{fakeRobotA, fakeRobotB},
		Marsupials: []simulations.Marsupial{
			simulations.NewMarsupial(fakeRobotA, fakeRobotB),
		},
	})

	assert.Equal(t, "cloudsim_sim.ign", cmd[0])
	assert.Equal(t, "worldName:=tunnel_circuit_practice_02", cmd[1])
	assert.Equal(t, fmt.Sprintf("durationSec:=%d", int(duration.Seconds())), cmd[2])
	assert.Equal(t, "headless:=true", cmd[3])
	assert.Equal(t, fmt.Sprintf("seed:=%d", seeds[runIndex]), cmd[4])
	assert.Equal(t, fmt.Sprintf("websocketAuthKey:=%s", token), cmd[5])
	assert.Equal(t, fmt.Sprintf("websocketAdminAuthKey:=%s", token), cmd[6])
	assert.Equal(t, fmt.Sprintf("websocketMaxConnections:=%d", maxConn), cmd[7])

	assert.Equal(t, fmt.Sprintf("robotName1:=%s", fakeRobotA.Name()), cmd[8])
	assert.Equal(t, fmt.Sprintf("robotConfig1:=%s", fakeRobotA.Kind()), cmd[9])

	assert.Equal(t, fmt.Sprintf("robotName2:=%s", fakeRobotB.Name()), cmd[10])
	assert.Equal(t, fmt.Sprintf("robotConfig2:=%s", fakeRobotB.Kind()), cmd[11])

	assert.Equal(t, fmt.Sprintf("marsupial1:=%s:%s", fakeRobotA.Name(), fakeRobotB.Name()), cmd[12])
}
