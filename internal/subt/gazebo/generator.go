package gazebo

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strings"
	"time"
)

const (
	keyDurationSec           = "durationSec"
	keyHeadless              = "headless"
	keySeed                  = "seed"
	keyWebsocketAuthKey      = "websocketAuthKey"
	keyWebsocketAdminAuthKey = "websocketAdminAuthKey"
	keyWebsocketMaxConn      = "websocketMaxConnections"
	keyRobotName             = "robotName"
	keyRobotConfig           = "robotConfig"
	keyMarsupial             = "marsupial"
)

// LaunchConfig includes the information to create the launch command arguments needed to launch gazebo server.
type LaunchConfig struct {
	Worlds                  string
	WorldMaxSimSeconds      time.Duration
	Seeds                   []int
	WorldIndex              *int
	RunIndex                *int
	AuthorizationToken      *string
	MaxWebsocketConnections int
	Robots                  []simulations.Robot
	Marsupials              []simulations.Marsupial
}

// Generate generates the needed arguments to initialize Gazebo server.
func Generate(params LaunchConfig) []string {
	worlds := ign.StrToSlice(params.Worlds)

	launchWorldName := worlds[0]
	if params.WorldIndex != nil {
		launchWorldName = worlds[*params.WorldIndex]
	}

	// We split by ";" (semicolon), in case the configured worldToLaunch string has arguments.
	// eg. 'tunnel_circuit_practice.ign;worldName:=tunnel_circuit_practice_01'
	var cmd []string
	cmd = strings.Split(launchWorldName, ";")

	// Set the simulation time limit
	cmd = append(cmd, fmt.Sprintf("%s:=%d", keyDurationSec, int(params.WorldMaxSimSeconds.Seconds())))

	// Set headless
	cmd = append(cmd, fmt.Sprintf("%s:=%s", keyHeadless, "true"))

	// Get the Seed for this run
	if len(params.Seeds) > 0 {
		seed := params.Seeds[0]
		if params.RunIndex != nil {
			seed = params.Seeds[*params.RunIndex]
		}

		cmd = append(cmd, fmt.Sprintf("%s:=%d", keySeed, seed))
	}

	// Set the authorization token if it exists
	if params.AuthorizationToken != nil {
		cmd = append(cmd, fmt.Sprintf("%s:=%s", keyWebsocketAuthKey, *params.AuthorizationToken))
		cmd = append(cmd, fmt.Sprintf("%s:=%s", keyWebsocketAdminAuthKey, *params.AuthorizationToken))
	}

	// Set the maximum number of websocket connections. This can be removed
	// when websocket scaling across multiple machines is implemented.
	cmd = append(cmd, fmt.Sprintf("%s:=%d", keyWebsocketMaxConn, params.MaxWebsocketConnections))

	// Pass Robot names and types to the gzserver Pod.
	// robotName1:=xxx robotConfig1:=yyy robotName2:=xxx robotConfig2:=yyy (Note the numbers).
	for i, robot := range params.Robots {
		cmd = append(cmd,
			fmt.Sprintf("%s%d:=%s", keyRobotName, i+1, robot.Name()),
			fmt.Sprintf("%s%d:=%s", keyRobotConfig, i+1, robot.Kind()),
		)
	}

	// Pass marsupial names to the gzserver Pod.
	// marsupialN:=<parent>:<child>
	for i, marsupial := range params.Marsupials {
		cmd = append(cmd, fmt.Sprintf("%s%d:=%s:%s", keyMarsupial, i+1, marsupial.Parent().Name(), marsupial.Child().Name()))
	}

	return cmd
}
