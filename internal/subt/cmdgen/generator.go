package cmdgen

import (
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"strings"
	"time"
)

var (
	// ErrEmptyWorld is returned when an empty world name is passed when calling CommsBridge.
	ErrEmptyWorld = errors.New("empty world")
	// ErrInvalidRobot is returned when an invalid robot is passed when calling CommsBridge.
	ErrInvalidRobot = errors.New("invalid robot")
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

// CommsBridgeConfig includes the information needed to generate the arguments for the
// comms bridge container.
type CommsBridgeConfig struct {
	// World represents a command to launch a certain world with the needed parameters.
	World string
	// RobotNumber is a number from 0 to the max number of robots - 1.
	// It usually is the index when looping over the list of robots for a certain simulation.
	RobotNumber int
	// Robot is the robot launched in the field computer linked to the comms bridge container.
	Robot simulations.Robot
	// ChildMarsupial is true if the robot given in this configuration is a marsupial child.
	ChildMarsupial bool
}

// CommsBridge generates the arguments needed to run in the comms bridge container.
func CommsBridge(config CommsBridgeConfig) ([]string, error) {
	params := strings.Split(config.World, ";")
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

	if config.Robot == nil {
		return nil, ErrInvalidRobot
	}

	return []string{
		worldNameParam,
		fmt.Sprintf("robotName%d:=%s", config.RobotNumber+1, config.Robot.GetName()),
		fmt.Sprintf("robotConfig%d:=%s", config.RobotNumber+1, config.Robot.GetKind()),
		"headless:=true",
		fmt.Sprintf("marsupial:=%t", config.ChildMarsupial),
	}, nil
}

// MapAnalysisConfig has the fields needed to configure a Mapping server container.
type MapAnalysisConfig struct {
	// World is a gazebo world with parameters.
	// Example:
	// 	"tunnel_circuit_practice.ign;worldName:=tunnel_circuit_practice_01"
	World string

	// Robots includes the information about all simulation robots.
	Robots []simulations.Robot
}

// MapAnalysis generates a set of arguments to configure the Mapping server container.
func MapAnalysis(config MapAnalysisConfig) ([]string, error) {
	params := strings.Split(config.World, ";")
	var worldName string
	for _, param := range params {
		if strings.Index(param, "worldName:=") != -1 {
			if _, err := fmt.Sscanf(param, "worldName:=%s", &worldName); err != nil {
				return nil, err
			}
			break
		}
	}

	if worldName == "" {
		return nil, ErrEmptyWorld
	}

	pdc := fmt.Sprintf("pdc:=%s.pdc", worldName)
	gt := fmt.Sprintf("gt:=%s.csv", worldName)

	robots := make([]string, len(config.Robots))
	for i, r := range config.Robots {
		robots[i] = fmt.Sprintf("robot:=%s", r.GetName())
	}

	return append([]string{pdc, gt}, robots...), nil
}
