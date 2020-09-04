package gazebo

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strconv"
	"strings"
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

// LaunchParams includes the information to create the launch command arguments needed to launch gazebo server.
type LaunchParams struct {
	Worlds                  string
	WorldMaxSimSeconds      string
	Seeds                   *string
	WorldIndex              *int
	RunIndex                *int
	AuthorizationToken      *string
	MaxWebsocketConnections int
	Robots                  []simulations.Robot
	Marsupials              []simulations.Marsupial
}

// GenerateLaunchArgs generates the needed arguments to initialize Gazebo server.
func GenerateLaunchArgs(params LaunchParams) ([]string, error) {
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
	cmd = append(cmd, fmt.Sprintf("%s:=%s", keyDurationSec, params.WorldMaxSimSeconds))

	// Set headless
	cmd = append(cmd, fmt.Sprintf("%s:=%s", keyHeadless, "true"))

	// Get the configured Seed for this run
	if params.Seeds != nil {
		seeds, err := SplitInt(*params.Seeds, ",")
		if err != nil {
			return nil, err
		}

		var seed int

		seed = seeds[0]
		if params.RunIndex != nil {
			seed = seeds[*params.RunIndex]
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
			fmt.Sprintf("%s%d:=%s", keyRobotConfig, i+1, robot.Type()),
		)
	}

	// Pass marsupial names to the gzserver Pod.
	// marsupialN:=<parent>:<child>
	for i, marsupial := range params.Marsupials {
		cmd = append(cmd, fmt.Sprintf("%s%d:=%s:%s", keyMarsupial, i+1, marsupial.Parent().Name(), marsupial.Child().Name()))
	}

	return cmd, nil
}

// SplitInt splits the given string by the given separator. It also parses all the content to int.
// It will return an error if the parse fails.
func SplitInt(s string, sep string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	noSpaces := strings.TrimSpace(s)
	noSpaces = strings.TrimPrefix(noSpaces, sep)
	noSpaces = strings.TrimSuffix(noSpaces, sep)
	var result []int
	for _, numStr := range strings.Split(noSpaces, ",") {
		numStr = strings.TrimSpace(numStr)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}
