package cmdgen

import (
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"strings"
	"time"
)

// ErrEmptyWorld is returned when an empty world name is passed when calling CommsBridge.
var ErrEmptyWorld = errors.New("empty world")

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
	keyRos                   = "ros"
)

// GazeboConfig includes the information to create the launch command arguments needed to launch gazebo server.
type GazeboConfig struct {
	// World is a gazebo world with parameters.
	// Example:
	// 	"tunnel_circuit_practice.ign;worldName:=tunnel_circuit_practice_01"
	World string

	// WorldMaxSimSeconds is the total amount of seconds that a simulation can run.
	WorldMaxSimSeconds time.Duration

	// Seed is used to randomly generate a world.
	// If no seed is provided, gazebo will generate its own seed.
	Seed *int

	// AuthorizationToken has the token used for the gazebo websocket server.
	AuthorizationToken *string

	// MaxWebsocketConnections determines the maximum amount of connections that can be established with
	// the websocket server.
	MaxWebsocketConnections int

	// Robots is a group of robots that will be used for this simulation.
	Robots []simulations.Robot

	// Marsupials is a group of parent-child pair robots. The robots used as parent and child should be in the
	// Robots slice as well.
	Marsupials []simulations.Marsupial

	// RosEnabled is used to enable ros when launching gazebo server.
	RosEnabled bool
}

// Gazebo generates the needed arguments to initialize the gzserver.
func Gazebo(params GazeboConfig) []string {
	launchWorldName := params.World

	// We split by ";" (semicolon), in case the configured worldToLaunch string has arguments.
	// eg. 'tunnel_circuit_practice.ign;worldName:=tunnel_circuit_practice_01'
	var cmd []string
	cmd = strings.Split(launchWorldName, ";")

	// Set the simulation time limit
	cmd = append(cmd, fmt.Sprintf("%s:=%d", keyDurationSec, int(params.WorldMaxSimSeconds.Seconds())))

	// Set headless
	cmd = append(cmd, fmt.Sprintf("%s:=%s", keyHeadless, "true"))

	// Set the Seed for this run
	if params.Seed != nil {
		cmd = append(cmd, fmt.Sprintf("%s:=%d", keySeed, *params.Seed))
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
			fmt.Sprintf("%s%d:=%s", keyRobotName, i+1, robot.GetName()),
			fmt.Sprintf("%s%d:=%s", keyRobotConfig, i+1, robot.GetKind()),
		)
	}

	// Pass marsupial names to the gzserver Pod.
	// marsupialN:=<parent>:<child>
	for i, marsupial := range params.Marsupials {
		cmd = append(cmd, fmt.Sprintf("%s%d:=%s:%s", keyMarsupial, i+1, marsupial.GetParent().GetName(), marsupial.GetChild().GetName()))
	}

	cmd = append(cmd, fmt.Sprintf("%s:=%t", keyRos, params.RosEnabled))

	return cmd
}

// CommsBridge generates the arguments needed to run in the comms bridge container.
func CommsBridge(world string, robotNumber int, robotName string, robotType string, childMarsupial bool) ([]string, error) {
	params := strings.Split(world, ";")
	var worldNameParam string
	for _, param := range params {
		if strings.Index(param, "worldName:=") != -1 {
			worldNameParam = param
			break
		}
	}

	if worldNameParam == "" {
		return nil, ErrEmptyWorld
	}

	return []string{
		worldNameParam,
		fmt.Sprintf("robotName%d:=%s", robotNumber+1, robotName),
		fmt.Sprintf("robotConfig%d:=%s", robotNumber+1, robotType),
		"headless:=true",
		fmt.Sprintf("marsupial:=%t", childMarsupial),
	}, nil
}
