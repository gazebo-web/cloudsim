package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// Integration tests for Simulation related routes.

type createSimulationTest struct {
	uriTest
	circuit    string
	owner      string
	robotName  *string
	expCreator string
}

type createSimulationCreditTest struct {
	uriTest
	circuit string
}

type getSimulationsTest struct {
	uriTest
	expSimNames []string
	circuit     *string
}

type getSimulationsMetadataTest struct {
	uriTest
	isAdmin    bool
	isMultisim bool
}

// createSimulationDeployment creates simulation deployments for testing.
func createSimulationDeployment(t *testing.T, ctx context.Context, db *gorm.DB, jwt *testJWT, simName string,
	circuit string, owner string, robotName *string, robotType *string) *sim.SimulationDeployment {

	// Prepare sim creation request
	uri := uriTest{
		testDesc:          "createSim",
		URL:               "/1.0/simulations",
		jwtGen:            jwt,
		expErrMsg:         nil,
		ignoreErrorBody:   false,
		ignoreOptionsCall: false,
	}
	if robotName == nil {
		robotName = sptr("X1")
	}
	if robotType == nil {
		robotType = sptr("X1_SENSOR_CONFIG_1")
	}
	createSubtForm := map[string]string{
		"name":        simName,
		"owner":       owner,
		"circuit":     circuit,
		"robot_name":  *robotName,
		"robot_type":  *robotType,
		"robot_image": "infrastructureascode/aws-cli:latest",
	}

	// Request sim creation
	simDep := &sim.SimulationDeployment{}
	invokeURITestMultipartPOST(t, uri, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
		require.NoError(t, json.Unmarshal(*bslice, simDep), "Unable to unmarshal response", string(*bslice))
	})

	return simDep
}

func TestSimulationsRoute(t *testing.T) {
	// General test setup
	setup()

	defaultJWT := getDefaultTestJWT()

	createSimURI := "/1.0/simulations"

	invalidURI := "/1.0/simulations_inv"

	nonexistentJWT := newJWT(createJWTForIdentity(t, "User3"))
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	teamBAdmin := newJWT(createJWTForIdentity(t, "TeamBAdmin"))
	teamBUser1 := newJWT(createJWTForIdentity(t, "TeamBUser1"))
	appTeamMember := newJWT(createJWTForIdentity(t, "subtMember"))

	unauth := ign.NewErrorMessage(ign.ErrorUnauthorized)

	var teamBSimGroupID string
	var teamASimGroupID string

	vStix := "Virtual Stix"
	tunnel := "Tunnel Circuit"

	// Create a custom max_submissions rule for Tunnel Circuit (set to 1)
	tcLimitRule := &sim.CircuitCustomRule{
		Owner:    nil,
		Circuit:  &tunnel,
		RuleType: sim.MaxSubmissions,
		Value:    "1",
	}
	globals.Server.Db.Create(&tcLimitRule)

	createSimsTestsData := []createSimulationTest{
		{uriTest{"createSim - invalid uri", invalidURI, nil, ign.NewErrorMessage(ign.ErrorNameNotFound), true, true}, vStix, "", nil, ""},
		{uriTest{"createSim - with no jwt", createSimURI, nil, unauth, true, false}, vStix, "", nil, ""},
		{uriTest{"createSim - valid invocation with jwt", createSimURI, teamAUser1, nil, false, false}, vStix, "TeamA", nil, "TeamAUser1"},
		{uriTest{"createSim - valid invocation for TeamB with default jwt", createSimURI, defaultJWT, nil, false, false}, vStix, "TeamB", nil, sysAdminForTest},
		{uriTest{"createSim - TeamBAdmin should not be able to create sim for TeamA", createSimURI, teamBAdmin, unauth, false, false}, vStix, "TeamA", nil, ""},
		{uriTest{"createSim - TeamBAdmin can create sim for TeamB", createSimURI, teamBAdmin, nil, false, false}, tunnel, "TeamB", nil, "TeamBAdmin"},
		{uriTest{"createSim - TeamBAdmin should not be able to create a 2nd sim for TeamB", createSimURI, teamBAdmin, sim.NewErrorMessage(sim.ErrorCircuitSubmissionLimitReached), false, false}, tunnel, "TeamB", nil, "TeamBAdmin"},
		{uriTest{"createSim - robotName longer than max", createSimURI, teamBAdmin, ign.NewErrorMessage(ign.ErrorFormInvalidValue), false, false}, tunnel, "TeamB", sptr("a25charsname1111111111111"), "TeamBAdmin"},
	}

	for i, test := range createSimsTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			rName := "robot1"
			if test.robotName != nil {
				rName = *test.robotName
			}
			createSubtForm := map[string]string{
				"name":        fmt.Sprintf("sim%d", i),
				"owner":       test.owner,
				"circuit":     test.circuit,
				"robot_name":  rName,
				"robot_type":  "X1_SENSOR_CONFIG_1",
				"robot_image": "infrastructureascode/aws-cli:latest",
			}
			invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				dep := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
				assert.Equal(t, test.owner, *dep.Owner)
				assert.Equal(t, test.expCreator, *dep.Creator)
				assert.True(t, *dep.Private)
				assert.Equal(t, "subt", *dep.Application)
				assert.Equal(t, 0, dep.MultiSim)
				if teamBSimGroupID == "" && *dep.Creator == "TeamBAdmin" {
					// save the created dep groupID
					// HACK
					teamBSimGroupID = *dep.GroupID
				}
				if teamASimGroupID == "" && *dep.Creator == "TeamAUser1" {
					// save the created dep groupID
					// HACK
					teamASimGroupID = *dep.GroupID
				}
			})
		})
	}

	uri := "/1.0/simulations"

	getSimulationsTestsData := []getSimulationsTest{
		{uriTest{"getSims - invalid uri", invalidURI, nil, ign.NewErrorMessage(ign.ErrorNameNotFound), true, true}, nil, nil}, {uriTest{"getSims - with no jwt", uri, nil, ign.NewErrorMessage(ign.ErrorUnauthorized), true, false}, nil, nil},
		{uriTest{"getSims - invocation with non existing user", uri, nonexistentJWT, ign.NewErrorMessage(ign.ErrorAuthNoUser), false, false}, nil, nil},
		{uriTest{"getSims - valid invocation with jwt", uri, defaultJWT, nil, false, false}, []string{"sim5", "sim3", "sim2"}, nil},
		{uriTest{"getSims - valid invocation with user valid jwt", uri, teamAUser1, nil, false, false}, []string{"sim2"}, nil},
		{uriTest{"getSims - valid invocation filtering by circuit with jwt", uri, defaultJWT, nil, false, false}, []string{"sim5"}, sptr(tunnel)},
		{uriTest{"getSims - valid invocation filtering by circuit with user valid jwt", uri, teamAUser1, nil, false, false}, []string{"sim2"}, sptr(vStix)},
		{uriTest{"getSims - valid invocation filtering by circuit with admin jwt", uri, defaultJWT, nil, false, false}, []string{"sim3", "sim2"}, sptr(vStix)},
	}

	for _, test := range getSimulationsTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			if test.circuit != nil {
				test.uriTest.URL = fmt.Sprintf("%s?circuit=%s", test.uriTest.URL, *test.circuit)
			}
			invokeURITest(t, test.uriTest, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				sims := []sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &sims), "Unable to unmarshal response", string(*bslice))
				for i, n := range test.expSimNames {
					assert.Equal(t, n, *sims[i].Name)
				}
			})
		})
	}
	teamBSim := uri + "/" + teamBSimGroupID
	getSingleSimTestsData := []getSimulationsTest{
		{uriTest{"getSingleSim - invalid uri", uri + "/invalid", defaultJWT, ign.NewErrorMessage(ign.ErrorNameNotFound), true, true}, nil, nil},
		{uriTest{"getSingleSim - with no jwt", teamBSim, nil, unauth, true, false}, nil, nil},
		{uriTest{"getSingleSim - admin from TeamB is OK", teamBSim, teamBAdmin, nil, false, false}, []string{"sim5"}, nil},
		{uriTest{"getSingleSim - user from TeamB is OK", teamBSim, teamBUser1, nil, false, false}, []string{"sim5"}, nil},
		{uriTest{"getSingleSim - user from TeamA cannot read", teamBSim, teamAUser1, unauth, false, false}, nil, nil},
	}

	for _, test := range getSingleSimTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			invokeURITest(t, test.uriTest, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				sim := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &sim), "Unable to unmarshal response", string(*bslice))
				assert.Equal(t, test.expSimNames[0], *sim.Name)
			})
		})
	}

	teamASim := uri + "/" + teamASimGroupID
	stopSimTestsData := []getSimulationsTest{
		{uriTest{"stopSimB - invalid uri", uri + "/invalid", defaultJWT, ign.NewErrorMessage(ign.ErrorNameNotFound), true, true}, nil, nil},
		{uriTest{"stopSimB - with no jwt", teamBSim, nil, unauth, true, false}, nil, nil},
		{uriTest{"stopSimB - user from TeamA cannot stop simB", teamBSim, teamAUser1, unauth, false, false}, nil, nil},
		{uriTest{"stopSimB - user from TeamB should not be able to stop a Tunnel Circuit", teamBSim, teamBUser1, unauth, false, false}, nil, nil},
		{uriTest{"stopSimB - App admin should be able to stop a Tunnel Circuit", teamBSim, appTeamMember, nil, false, false}, []string{"sim5"}, nil},
		{uriTest{"stopSimA - App admin can stop sim OK", teamASim, appTeamMember, nil, false, false}, []string{"sim2"}, nil},
	}
	for _, test := range stopSimTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			invokeURITestWithArgs(t, test.uriTest, "DELETE", nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				dep := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))
				assert.Equal(t, test.expSimNames[0], *dep.Name)
			})
		})
	}
}

func TestGetSimExtra(t *testing.T) {
	setup()
	db := globals.Server.Db

	const uri = "/1.0/simulations"

	sysAdmin := getDefaultTestJWT()
	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	singleSimCircuit := "Virtual Stix"
	multiSimCircuit := "Urban Practice 1"

	ctx := context.Background()

	singleSim := createSimulationDeployment(
		t, ctx, db, teamAUser1, "TestSingleGroupIDSimulation", singleSimCircuit, "TeamA", nil, nil,
	)
	multiSim := createSimulationDeployment(
		t, ctx, db, teamAUser1, "TestMultiGroupIDSimulation", multiSimCircuit, "TeamA", nil, nil,
	)

	simURL := fmt.Sprintf("%s/%s", uri, *singleSim.GroupID)
	multisimURL := fmt.Sprintf("%s/%s", uri, *multiSim.GroupID)

	getSingleSimMetadataTests := []getSimulationsMetadataTest{
		{uriTest{"getSimMetadata - user does not get metadata", simURL, teamAUser1, nil, true, false}, false, false},
		{uriTest{"getSimMetadata - sysadmin does not get metadata", simURL, sysAdmin, nil, true, false}, true, false},
		{uriTest{"getMultiSimMetadata - user does not get metadata", multisimURL, teamAUser1, nil, true, false}, false, true},
		{uriTest{"getMultiSimMetadata - sysadmin does get metadata", multisimURL, sysAdmin, nil, true, false}, true, true},
	}

	for _, test := range getSingleSimMetadataTests {
		t.Run(test.testDesc, func(t *testing.T) {
			invokeURITest(t, test.uriTest, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				dep := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &dep), "Unable to unmarshal response", string(*bslice))

				extra, _ := sim.ReadExtraInfoSubT(&dep)

				assert.Nil(t, extra.WorldIndex)
				assert.Nil(t, extra.RunIndex)
			})
		})
	}
}

func TestCheckCredit(t *testing.T) {
	setup()

	teamBAdmin := newJWT(createJWTForIdentity(t, "TeamBAdmin"))

	createSimURI := "/1.0/simulations"

	vStix := "Virtual Stix"
	tunnel := "Tunnel Circuit"
	practice := "Tunnel Practice 1"

	checkSimulationCreditTestData := []createSimulationCreditTest{
		{uriTest{"checkSimCredits - Robot credit sum is less than credits limit", createSimURI, teamBAdmin, nil, false, false}, vStix},
		{uriTest{"checkSimCredits - Robot credit sum is greater than credits limit", createSimURI, teamBAdmin, sim.NewErrorMessage(sim.ErrorCreditsExceeded), false, false}, practice},
		{uriTest{"checkSimCredits - Null/Infinite credits limit", createSimURI, teamBAdmin, nil, false, false}, tunnel},
	}

	for i, test := range checkSimulationCreditTestData {
		t.Run(test.testDesc, func(t *testing.T) {

			rName := "robot1"

			createSubtForm := map[string]string{
				"name":        fmt.Sprintf("sim%d", i),
				"owner":       "TeamB",
				"circuit":     test.circuit,
				"robot_name":  rName,
				"robot_type":  "X1_SENSOR_CONFIG_1",
				"robot_image": "infrastructureascode/aws-cli:latest",
			}

			invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				sim := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &sim), "Unable to unmarshal response", string(*bslice))
			})
		})
	}
}

func TestGetRemainingSubmissionsRoute(t *testing.T) {
	// General test setup
	setup()

	// Route parameters
	circuit := "Tunnel Circuit"
	owner := sysAdminForTest
	URI := fmt.Sprintf("/1.0/%s/remaining_submissions/%s", circuit, owner)
	testRoute := uriTest{
		"Get remaining submissions",
		URI,
		getDefaultTestJWT(),
		nil,
		true,
		true,
	}

	// Create special rules
	type remainingSubmissionsTest struct {
		rule                         *sim.CircuitCustomRule
		expectedRemainingSubmissions *int
	}
	customRules := []remainingSubmissionsTest{
		{
			nil,
			nil,
		},
		{
			&sim.CircuitCustomRule{
				Owner:    nil,
				Circuit:  nil,
				RuleType: sim.MaxSubmissions,
				Value:    "11",
			},
			intptr(9),
		},
		{
			&sim.CircuitCustomRule{
				Owner:    nil,
				Circuit:  sptr(circuit),
				RuleType: sim.MaxSubmissions,
				Value:    "22",
			},
			intptr(19),
		},
		{
			&sim.CircuitCustomRule{
				Owner:    sptr(owner),
				Circuit:  nil,
				RuleType: sim.MaxSubmissions,
				Value:    "33",
			},
			intptr(29),
		},
		{
			&sim.CircuitCustomRule{
				Owner:    sptr(owner),
				Circuit:  sptr(circuit),
				RuleType: sim.MaxSubmissions,
				Value:    "44",
			},
			intptr(39),
		},
	}

	// Custom rules with ascending priority are introduced in order. A fake
	// simulation is created before each rule insertion to increase the
	// owner's number of submissions in order to check the number decreases.
	for i, test := range customRules {
		t.Run(testRoute.testDesc+string(i), func(t *testing.T) {
			// Create a fake simulation
			globals.Server.Db.Create(&sim.SimulationDeployment{
				Owner:         sptr(owner),
				GroupID:       sptr(uuid.NewV4().String()),
				ExtraSelector: sptr(circuit),
			})

			// Create the next custom rules
			if test.rule != nil {
				globals.Server.Db.Create(&test.rule)
			}

			invokeURITest(t, testRoute, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				response := sim.RemainingSubmissions{}
				count := 0
				require.NoError(t, json.Unmarshal(*bslice, &response), "Unable to unmarshal response", string(*bslice))
				globals.Server.Db.Model(&sim.SimulationDeployment{}).Count(&count)

				if test.expectedRemainingSubmissions == nil {
					assert.Equal(
						t,
						test.expectedRemainingSubmissions,
						response.RemainingSubmissions,
					)
				} else {
					assert.Equal(
						t,
						*test.expectedRemainingSubmissions,
						*response.RemainingSubmissions,
					)
				}
			})
		})
	}
}

func TestSetCustomRuleRoute(t *testing.T) {
	// General test setup
	setup()

	// User setup
	sysAdmin := getDefaultTestJWT()
	subtAdmin := newJWT(createJWTForIdentity(t, "subtAdmin"))
	teamUser := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	// Route parameters
	circuit := "Tunnel Circuit"
	owner := sysAdminForTest
	URI := fmt.Sprintf("/1.0/rules")

	// Create new rules
	type customRulesTest struct {
		uriTest  uriTest
		ruleType sim.CustomRuleType
		circuit  string
		owner    string
		value    string
		delete   bool
	}
	tests := []customRulesTest{
		{
			uriTest{
				"Create new rule for sysAdmin",
				URI,
				sysAdmin,
				nil,
				true,
				true,
			},
			sim.MaxSubmissions,
			circuit,
			owner,
			"123",
			false,
		},
		{
			uriTest{
				"Update rule for sysAdmin",
				URI,
				sysAdmin,
				nil,
				true,
				true,
			},
			sim.MaxSubmissions,
			circuit,
			owner,
			"777",
			true,
		},
		{
			uriTest{
				"Create new invalid rule for sysAdmin",
				URI,
				sysAdmin,
				ign.NewErrorMessage(ign.ErrorFormInvalidValue),
				true,
				true,
			},
			"invalid",
			circuit,
			owner,
			"123",
			false,
		},
		{
			uriTest{
				"Create new rule with team admin account",
				URI,
				subtAdmin,
				nil,
				true,
				true,
			},
			sim.MaxSubmissions,
			circuit,
			"TeamA",
			"10",
			false,
		},
		{
			uriTest{
				"Create new rule for another owner with admin account",
				URI,
				sysAdmin,
				nil,
				true,
				true,
			},
			sim.MaxSubmissions,
			circuit,
			"TeamAUser1",
			"10",
			false,
		},
		// TODO This is currently not supported. Admin privileges override invalid owner errors.
		//{
		//	uriTest{
		//		"Create new rule for invalid owner with sysadmin account",
		//		URI,
		//		sysAdmin,
		//		ign.NewErrorMessage(ign.ErrorFormInvalidValue),
		//		true,
		//		true,
		//	},
		//	sim.MaxSubmissions,
		//	circuit,
		//	"invalid",
		//	"10",
		//	false,
		//},
		{
			uriTest{
				"Create new rule for invalid circuit",
				URI,
				sysAdmin,
				ign.NewErrorMessage(ign.ErrorFormInvalidValue),
				true,
				true,
			},
			sim.MaxSubmissions,
			"invalid",
			owner,
			"10",
			false,
		},
		{
			uriTest{
				"Create new rule with regular team user",
				URI,
				teamUser,
				ign.NewErrorMessage(ign.ErrorUnauthorized),
				true,
				true,
			},
			sim.MaxSubmissions,
			circuit,
			"TeamA",
			"10",
			false,
		},
	}

	// Create/Update rules
	for _, test := range tests {
		t.Run(test.uriTest.testDesc, func(t *testing.T) {
			// Prepare the URL
			test.uriTest.URL = fmt.Sprintf("%s/%s/%s/%s/%s", test.uriTest.URL, test.circuit, test.owner,
				test.ruleType, test.value)

			invokeURITestWithArgs(t, test.uriTest, "PUT", nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				response := sim.CircuitCustomRule{}
				require.NoError(t, json.Unmarshal(*bslice, &response), "Unable to unmarshal response", string(*bslice))
				rule, err := sim.GetCircuitCustomRule(globals.Server.Db, circuit, test.owner, test.ruleType)
				require.NoError(t, err)
				assert.Equal(t, test.owner, *rule.Owner)
				assert.Equal(t, test.ruleType, rule.RuleType)
				assert.Equal(t, test.value, rule.Value)
				assert.Equal(t, test.owner, *response.Owner)
				assert.Equal(t, test.ruleType, response.RuleType)
				assert.Equal(t, test.value, response.Value)
			})
		})
	}

	// Delete rules
	for _, test := range tests {
		if !test.delete {
			continue
		}
		t.Run("Delete_"+test.uriTest.testDesc, func(t *testing.T) {
			// Prepare the URL
			test.uriTest.URL = fmt.Sprintf("%s/%s/%s/%s", URI, test.circuit, test.owner, test.ruleType)

			invokeURITestWithArgs(t, test.uriTest, "DELETE", nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				response := sim.CircuitCustomRule{}
				require.NoError(t, json.Unmarshal(*bslice, &response), "Unable to unmarshal response", string(*bslice))
				_, err := sim.GetCircuitCustomRule(globals.Server.Db, circuit, test.owner, test.ruleType)
				require.Error(t, err)
			})
		})
	}
}

func TestDownloadLogsRouter(t *testing.T) {
	// General test setup
	setup()

	// Create simulation deployments
	db := globals.Server.Db
	createSimulationDeployment := func(db *gorm.DB, owner string, groupID string, name string,
		multiSim int) sim.SimulationDeployment {
		extra := `{"circuit":"Tunnel Test 1","robots":[
			{"Name":"X1","Type":"X1_SENSOR_CONFIG_1","Image":"image"},
			{"Name":"X2","Type":"X2_SENSOR_CONFIG_1","Image":"image"}
		]}`
		simDep := sim.SimulationDeployment{
			ValidFor:         sptr("6h0m0s"),
			Owner:            sptr(owner),
			Creator:          sptr("test_user"),
			Private:          boolptr(true),
			GroupID:          sptr(groupID),
			DeploymentStatus: intptr(90),
			Platform:         sptr("subt"),
			Application:      sptr("subt"),
			Name:             sptr(name),
			MultiSim:         multiSim,
			Extra:            sptr(extra),
			ExtraSelector:    sptr("Tunnel Test 1"),
			StopOnEnd:        boolptr(false),
			Robots:           sptr("X1,X2"),
		}
		db.Create(&simDep)

		return simDep
	}
	simDepSingle := createSimulationDeployment(
		db, "TeamA", "test-single-groupID-simulation", "TestSingleSimSimulation", 0,
	)
	simDepMulti := createSimulationDeployment(
		db, "TeamA", "test-multi-groupID-simulation", "TestMultiSimSimulation", 1,
	)

	// User setup
	sysAdmin := getDefaultTestJWT()
	//teamUser := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	// Route parameters
	URI := "/1.0/simulations/%s/logs/file"

	// Setup tests
	robots := []string{"x1", "x2"}
	type logFileDownloadTest struct {
		uriTest  uriTest
		link     bool
		groupID  string
		filename string
		robot    *string
	}
	tests := []logFileDownloadTest{
		{
			uriTest{
				"Gazebo logs for single sim",
				fmt.Sprintf(URI, *simDepSingle.GroupID),
				sysAdmin,
				nil,
				true,
				true,
			},
			true,
			*simDepSingle.GroupID,
			fmt.Sprintf("%s.tar.gz", *simDepSingle.GroupID),
			nil,
		},
		{
			uriTest{
				"ROS logs for first robot of single sim",
				fmt.Sprintf(URI, *simDepSingle.GroupID),
				sysAdmin,
				nil,
				true,
				true,
			},
			true,
			*simDepSingle.GroupID,
			fmt.Sprintf("%s-fc-%s-commsbridge.tar.gz", *simDepSingle.GroupID, robots[0]),
			&robots[0],
		},
		{
			uriTest{
				"ROS logs for second robot of single sim",
				fmt.Sprintf(URI, *simDepSingle.GroupID),
				sysAdmin,
				nil,
				true,
				true,
			},
			true,
			*simDepSingle.GroupID,
			fmt.Sprintf("%s-fc-%s-commsbridge.tar.gz", *simDepSingle.GroupID, robots[1]),
			&robots[1],
		},
		{
			uriTest{
				"Summary for multisim",
				fmt.Sprintf(URI, *simDepMulti.GroupID),
				sysAdmin,
				nil,
				true,
				true,
			},
			true,
			*simDepMulti.GroupID,
			fmt.Sprintf("summary.json"),
			nil,
		},
		{
			uriTest{
				"Summary for multisim with robot parameter",
				fmt.Sprintf(URI, *simDepMulti.GroupID),
				sysAdmin,
				nil,
				true,
				true,
			},
			true,
			*simDepMulti.GroupID,
			fmt.Sprintf("summary.json"),
			&robots[0],
		},
	}
	// Create/Update rules
	for _, test := range tests {
		// Set query params
		params := make([]string, 0)
		if test.link {
			params = append(params, "link=true")
		}
		if test.robot != nil {
			params = append(params, fmt.Sprintf("robot=%s", *test.robot))
		}
		if len(params) > 0 {
			test.uriTest.URL += "?" + strings.Join(params, "&")
		}
		t.Run(test.uriTest.testDesc, func(t *testing.T) {
			invokeURITestWithArgs(t, test.uriTest, "GET", nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				var response string
				require.NoError(t, json.Unmarshal(*bslice, &response), "Unable to unmarshal response", string(*bslice))
				assert.Contains(t, response, test.filename)
			})
		})
	}
}

func TestGetCompetitionRobots(t *testing.T) {
	setup()

	// User setup
	sysAdmin := getDefaultTestJWT()
	subtAdmin := newJWT(createJWTForIdentity(t, "subtAdmin"))
	teamUser := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	URI := "/1.0/competition/robots"

	tests := []uriTest{
		{
			"Get circuit robots (sysadmin)",
			URI,
			sysAdmin,
			nil,
			true,
			true,
		},
		{
			"Get circuit robots (subt admin)",
			URI,
			subtAdmin,
			nil,
			true,
			true,
		},
		{
			"Get circuit robots (user)",
			URI,
			teamUser,
			nil,
			true,
			true,
		},
		{
			"Get circuit robots (anonymous)",
			URI,
			nil,
			ign.NewErrorMessage(ign.ErrorUnauthorized),
			true,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.testDesc, func(t *testing.T) {
			invokeURITestWithArgs(t, test, "GET", nil, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK callback
				var response map[string]sim.SubTRobotType
				require.NoError(t, json.Unmarshal(*bslice, &response), "Unable to unmarshal response", string(*bslice))
				for key := range sim.SubTRobotTypes {
					assert.Contains(t, response, key)
				}
			})
		})
	}
}
