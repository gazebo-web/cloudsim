package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntest "gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
	"testing"
	"time"
)

// Tests in this file must be run using a mocked EC2 client. In order to do this,
// all tests should call `useEC2NodeManager` after `setup` to ensure that the
// tests use a mocked EC2 client as the node manager. useEC2NodeManager
// returns a function that restores the project's node manager to its original
// value. The returned function should be deferred immediately after calling
// useEC2NodeManager.

// TestSuccessfulSimulations tests that a simulation is ran from start to finish.
func TestSuccessfulSimulations(t *testing.T) {
	// General test setup
	setup()
	// Set the EC2 Mock service as the node manager
	restoreHostSvc := useEC2NodeManager()
	defer restoreHostSvc()

	defaultJWT := getDefaultTestJWT()
	createSimURI := "/1.0/simulations"

	invalidURI := "/1.0/simulations_inv"

	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))
	teamBAdmin := newJWT(createJWTForIdentity(t, "TeamBAdmin"))

	unauth := ign.NewErrorMessage(ign.ErrorUnauthorized)

	var teamBSimGroupID string
	var teamASimGroupID string
	vStix := "Virtual Stix"
	tunnel := "Tunnel Circuit"
	createSimsTestsData := []createSimulationTest{
		{uriTest{"createSim - invalid uri", invalidURI, nil, ign.NewErrorMessage(ign.ErrorNameNotFound), true, true}, vStix, "", nil, ""},
		{uriTest{"createSim - with no jwt", createSimURI, nil, unauth, true, false}, vStix, "", nil, ""},
		{uriTest{"createSim - valid invocation with jwt", createSimURI, teamAUser1, nil, false, false}, vStix, "TeamA", nil, "TeamAUser1"},
		{uriTest{"createSim - valid invocation for TeamB with default jwt", createSimURI, defaultJWT, nil, false, false}, vStix, "TeamB", nil, sysAdminForTest},
		{uriTest{"createSim - TeamBAdmin should not be able to create sim for TeamA", createSimURI, teamBAdmin, unauth, false, false}, vStix, "TeamA", nil, ""},
		{uriTest{"createSim - TeamBAdmin can create sim for TeamB", createSimURI, teamBAdmin, nil, false, false}, tunnel, "TeamB", nil, "TeamBAdmin"},
		{uriTest{"createSim - robotName longer than max", createSimURI, teamBAdmin, ign.NewErrorMessage(ign.ErrorFormInvalidValue), false, false}, tunnel, "TeamB", sptr("a25charsname1111111111111"), "TeamBAdmin"},
	}

	for i, test := range createSimsTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			// Mock EC2 responses
			ec2Impl := NewEC2MockSuccessfulLaunch()
			sim.AssertMockedEC2(globals.EC2Svc).SetImpl(ec2Impl)

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

			t.Log("Running test ", test.testDesc)

			// WaitGroup is used to run tests sequentially
			var wg sync.WaitGroup
			if test.expErrMsg == nil {
				t.Log("WG add")
				wg.Add(1)
			}
			sim.SimServImpl.(*sim.Service).SetPoolEventsListener(func(poolEvent sim.PoolEvent, groupID string,
				result interface{}, em *ign.ErrMsg) {
				t.Log("WG done")
				wg.Done()
			})

			invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {
				// Status OK Callback
				simDep := sim.SimulationDeployment{}
				require.NoError(t, json.Unmarshal(*bslice, &simDep), "Unable to unmarshal response", string(*bslice))
				assert.Equal(t, test.owner, *simDep.Owner)
				assert.Equal(t, test.expCreator, *simDep.Creator)
				assert.True(t, *simDep.Private)
				assert.Equal(t, "subt", *simDep.Application)
				assert.Equal(t, 0, simDep.MultiSim)
				if teamBSimGroupID == "" && *simDep.Creator == "TeamBAdmin" {
					// save the created simDep groupID
					// HACK
					teamBSimGroupID = *simDep.GroupID
				}
				if teamASimGroupID == "" && *simDep.Creator == "TeamAUser1" {
					// save the created simDep groupID
					// HACK
					teamASimGroupID = *simDep.GroupID
				}
			})
			if test.expErrMsg == nil {
				t.Log("WG Wait")
				wg.Wait()
			}
		})
	}
}

// TestInsufficientCapacityRequeue tests that Cloudsim retries starting a simulation when not enough EC2 instances are
// available for a simulation request.
func TestInsufficientCapacityRequeue(t *testing.T) {
	// General test setup
	setup()
	// Set the EC2 Mock service as the node manager
	restoreHostSvc := useEC2NodeManager()
	defer restoreHostSvc()

	createSimURI := "/1.0/simulations"

	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))
	teamBAdmin := newJWT(createJWTForIdentity(t, "TeamBAdmin"))

	vStix := "Virtual Stix"
	tunnel := "Tunnel Circuit"
	createSimsTestsData := []createSimulationTest{
		{uriTest{"createSim - valid invocation with jwt", createSimURI, teamAUser1, nil, false, false}, vStix, "TeamA", nil, "TeamAUser1"},
		{uriTest{"createSim - TeamBAdmin can create sim for TeamB", createSimURI, teamBAdmin, nil, false, false}, tunnel, "TeamB", nil, "TeamBAdmin"},
	}

	for i, test := range createSimsTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			// Mock EC2 responses
			ec2Impl := NewEC2Mock()

			// Mock WaitUnitlInstanceStatusOk
			ec2Impl.SetMockFunction(Ec2OpWaitUntilInstanceStatusOk, FixedValues, false, nil)

			// Mock RunInstances
			mockValues := make([]interface{}, 0)
			// Set FixedValues to strict mode
			mockValues = append(mockValues, true)
			// Simulate that there's not enough resources 10 times
			for try := 0; try < 10; try++ {
				for j := 0; j < 6; j++ {
					mockValues = append(mockValues, ec2Impl.NewAWSErr(sim.AWSErrCodeInsufficientInstanceCapacity))
				}
			}
			// Simulate that there's not enough resources, then grant the resources
			mockValues = append(mockValues,
				ec2Impl.NewAWSErr(sim.AWSErrCodeInsufficientInstanceCapacity),
				ec2Impl.NewAWSErr(sim.AWSErrCodeInsufficientInstanceCapacity),
				// Check for available machines
				ec2Impl.NewAWSErr(sim.AWSErrCodeDryRunOperation),
				// EC2 Instance 1
				ec2Impl.NewAWSErr(sim.AWSErrCodeDryRunOperation),
				ec2Impl.NewReservation(fmt.Sprintf("i-test-1-%s", uuid.NewV4().String())),
				// EC2 Instance 2
				ec2Impl.NewAWSErr(sim.AWSErrCodeDryRunOperation),
				ec2Impl.NewReservation(fmt.Sprintf("i-test-2-%s", uuid.NewV4().String())),
			)
			ec2Impl.SetMockFunction(Ec2OpRunInstances, FixedValues, mockValues...)

			sim.AssertMockedEC2(globals.EC2Svc).SetImpl(ec2Impl)

			// Attempt to request instances up to 6 times
			sim.MaxAWSRetries = 6

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

			t.Log("Running test ", test.testDesc)

			// A WaitGroup is used to run tests sequentially
			var wg sync.WaitGroup
			if test.expErrMsg == nil {
				wg.Add(1)
			}
			//The pool listener acts as the tester as it is the last thing called when returning. Once it ends,
			//the work group is released to let other tests run.
			sim.SimServImpl.(*sim.Service).SetPoolEventsListener(func(poolEvent sim.PoolEvent, groupID string,
				result interface{}, em *ign.ErrMsg) {
				// If there are not enough resources, then the simulation will have been requeued and the lock should be
				// kept. Otherwise, this simulation is done and the lock can be released.
				if em == nil || em.ErrCode != ign.ErrorLaunchingCloudInstanceNotEnoughResources {
					defer wg.Done()
					if !ec2Impl.ValidateMockFunctions() {
						t.Fail()
					}
				}
			})

			invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {})

			//Lock until the current test finishes
			if test.expErrMsg == nil {
				wg.Wait()
			}
		})
	}
}

// TestFailedPodCreation tests a simulation that fails during simulation launch due to a pod never reaching the
// expected status.
func TestFailedPodCreation(t *testing.T) {
	// General test setup
	setup()
	// Set the EC2 Mock service as the node manager
	restoreHostSvc := useEC2NodeManager()
	defer restoreHostSvc()
	// Mock specific functions
	fixedValuesFuncWFMPC, _ := FixedValues(
		true,
		nil,
		nil,
		errors.New("mock"),
	)
	restoreWFMPC := ReplaceValue(&sim.WaitForMatchPodsCondition, func(ctx context.Context,
		c kubernetes.Interface, namespace string, opts metav1.ListOptions, condStr string,
		timeout time.Duration, condition sim.PodCondition) error {

		if value, ok := fixedValuesFuncWFMPC().(error); ok {
			return value
		}
		return nil
	})
	defer restoreWFMPC()

	createSimURI := "/1.0/simulations"

	teamAUser1 := newJWT(createJWTForIdentity(t, "TeamAUser1"))

	vStix := "Virtual Stix"
	createSimsTestsData := []createSimulationTest{
		{uriTest{"failedPodCreation - valid invocation with jwt", createSimURI, teamAUser1, nil, false, false}, vStix, "TeamA", nil, "TeamAUser1"},
	}

	for i, test := range createSimsTestsData {
		t.Run(test.testDesc, func(t *testing.T) {
			// Mock EC2 responses
			ec2Impl := NewEC2MockSuccessfulLaunch()
			ec2Impl.SetMockFunction(Ec2OpTerminateInstances, FixedValues, true,
				ec2Impl.NewAWSErr(sim.AWSErrCodeDryRunOperation),
				nil,
			)
			sim.AssertMockedEC2(globals.EC2Svc).SetImpl(ec2Impl)

			// Mock K8 responses
			k8Impl := NewK8Mock(context.Background())
			k8Impl.SetFixedMutatorsForPodCreation(
				MakePodStatusMutator(MakePodReadyStatus()),
				MakePodStatusMutator(MakePodFailedStatus()),
				MakePodStatusMutator(MakePodReadyStatus()),
			)
			sim.AssertMockedClientset(globals.KClientset).SetImpl(k8Impl)

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

			t.Log("Running test ", test.testDesc)

			// WaitGroup is used to run tests sequentially
			var wg sync.WaitGroup
			if test.expErrMsg == nil {
				t.Log("WG add")
				wg.Add(1)
			}
			sim.SimServImpl.(*sim.Service).SetPoolEventsListener(func(poolEvent sim.PoolEvent, groupID string,
				result interface{}, em *ign.ErrMsg) {
				// If there are not enough resources, then the simulation will have been requeued and the lock should be
				// kept. Otherwise, this simulation is done and the lock can be released.
				if poolEvent == sim.PoolRollbackFailedLaunch &&
					(em == nil || em.ErrCode != ign.ErrorLaunchingCloudInstanceNotEnoughResources) {
					defer wg.Done()
					if !ec2Impl.ValidateMockFunctions() {
						t.Fail()
					}
				}
			})

			invokeURITestMultipartPOST(t, test.uriTest, createSubtForm, func(bslice *[]byte, resp *igntest.AssertResponse) {})

			if test.expErrMsg == nil {
				t.Log("WG Wait")
				wg.Wait()
			}
		})
	}
}
