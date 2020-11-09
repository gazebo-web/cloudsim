package simulations

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

////////////////////
// Simulation error codes
////////////////////

// SimErrCode is the type for all simulations related error codes
type SimErrCode int64

const (
	// ErrorMarkingLocalNodeAsFree is triggered when there is an error updating a node
	// labels to mark a minikube node as free.
	ErrorMarkingLocalNodeAsFree SimErrCode = 5501
	// ErrorMarkingLocalNodeAsUsed is triggered when there is an error updating a node
	// labels to mark a minikube node as being used.
	ErrorMarkingLocalNodeAsUsed SimErrCode = 5502
	// ErrorFreeNodeNotFound is triggered when there is no free node to use (minikube)
	ErrorFreeNodeNotFound SimErrCode = 5503
	// ErrorLabeledNodeNotFound is triggered when a node cannot be found using a label selector.
	ErrorLabeledNodeNotFound SimErrCode = 5504
	// ErrorCreatingRunningSimulationNode is triggered when a RunningSimulation cannot be created.
	ErrorCreatingRunningSimulationNode SimErrCode = 5505
	// ErrorOwnerSimulationsLimitReached is triggered when an Owner wants to launch in parallel
	// more than the allowed simulations.
	ErrorOwnerSimulationsLimitReached SimErrCode = 5506
	// ErrorCircuitSubmissionLimitReached is triggered when an Owner wants to
	// launch a simulation to a circuit that has a limit on the number of
	// submissions.
	ErrorCircuitSubmissionLimitReached SimErrCode = 5507
	// ErrorRuleForOwnerNotFound is triggered when an Admin attempts to delete
	// a rule that does not exist for a specific Owner
	ErrorRuleForOwnerNotFound SimErrCode = 5508
	// ErrorRuleNotFound is triggered when a nonexistent rule is set/deleted
	ErrorRuleNotFound SimErrCode = 5509
	// ErrorInvalidScore is triggered when the score file for a simulation cannot be accessed
	ErrorInvalidScore SimErrCode = 5510
	// ErrorInvalidSummary is triggered when the simulation summary file for a simulation cannot be accessed
	ErrorInvalidSummary SimErrCode = 5511
	// ErrorFailedToUploadLogs is triggered when simulations logs fail to upload
	ErrorFailedToUploadLogs SimErrCode = 5512
	// ErrorCreditsExceeded is triggered when the robot credit sum exceeds the circuit credits limit.
	ErrorCreditsExceeded SimErrCode = 5513
	// ErrorCircuitRuleNotFound is triggered when the circuit rule doesn't exist.
	ErrorCircuitRuleNotFound SimErrCode = 5514
	// ErrorFailedToGetLiveLogs is triggered when simulation live logs cannot be get
	ErrorFailedToGetLiveLogs SimErrCode = 5515
	// ErrorRobotIdentifierNotFound is triggered when a robot identifier cannot be found.
	ErrorRobotIdentifierNotFound SimErrCode = 5516
	// ErrorCompetitionNotStarted is triggered when a competition hasn't started yet.
	ErrorCompetitionNotStarted SimErrCode = 5517
	// ErrorNotQualified is triggered when a user/team is not qualified to compete on a circuit.
	ErrorNotQualified SimErrCode = 5518
	// ErrorLaunchHeldSimulation is triggered when an error is found while launching a held simulation.
	ErrorLaunchHeldSimulation SimErrCode = 5519
	// ErrorInvalidRobotImage is triggered when an owner attempts to start a simulation with a robot image that does
	// not belong to them
	ErrorInvalidRobotImage SimErrCode = 5520
	// ErrorInvalidMarsupialSpecification is triggered when an invalid SubT marsupial pair is specified.
	ErrorInvalidMarsupialSpecification SimErrCode = 5521
	// ErrorRobotModelLimitReached is triggered when attempting to launch a simulation with too many robots of a
	// single model.
	ErrorRobotModelLimitReached SimErrCode = 5522
	// ErrorSubmissionDeadlineReached is triggered when a submission deadline for a certain circuit has been reached.
	ErrorSubmissionDeadlineReached SimErrCode = 5523
	// ErrorLaunchSupersededSimulation is triggered when a held simulation is superseded.
	ErrorLaunchSupersededSimulation SimErrCode = 5224
)

// NewErrorMessageWithBase receives an error code and a root error
// and returns a pointer to an ErrMsg.
func NewErrorMessageWithBase(err SimErrCode, base error) *ign.ErrMsg {
	em := NewErrorMessage(err)
	em.BaseError = ign.WithStack(base)
	return em
}

// NewErrorMessage is a convenience function that receives an error code
// and returns a pointer to an ErrMsg.
func NewErrorMessage(err SimErrCode) *ign.ErrMsg {
	em := ErrorMessage(err)
	return &em
}

// ErrorMessage receives an error code and generate an error message response
func ErrorMessage(err SimErrCode) ign.ErrMsg {

	em := ign.ErrorMessageOK()

	em.ErrID = uuid.NewV4().String()

	switch err {
	case ErrorMarkingLocalNodeAsFree:
		em.Msg = "Error marking minikube node as free."
		em.ErrCode = int(ErrorMarkingLocalNodeAsFree)
		em.StatusCode = http.StatusInternalServerError
	case ErrorMarkingLocalNodeAsUsed:
		em.Msg = "Error marking minikube node as being used."
		em.ErrCode = int(ErrorMarkingLocalNodeAsUsed)
		em.StatusCode = http.StatusInternalServerError
	case ErrorFreeNodeNotFound:
		em.Msg = "There are no free minikubes nodes to use."
		em.ErrCode = int(ErrorFreeNodeNotFound)
		em.StatusCode = http.StatusInternalServerError
	case ErrorLabeledNodeNotFound:
		em.Msg = "Node could not be found using a label selector."
		em.ErrCode = int(ErrorLabeledNodeNotFound)
		em.StatusCode = http.StatusInternalServerError
	case ErrorCreatingRunningSimulationNode:
		em.Msg = "RunningSimulation node could not be created."
		em.ErrCode = int(ErrorCreatingRunningSimulationNode)
		em.StatusCode = http.StatusInternalServerError
	case ErrorOwnerSimulationsLimitReached:
		em.Msg = "Simultaneous simulations limit reached."
		em.ErrCode = int(ErrorOwnerSimulationsLimitReached)
		em.StatusCode = http.StatusBadRequest
	case ErrorCircuitSubmissionLimitReached:
		em.Msg = "Circuit simulation submission limit reached."
		em.ErrCode = int(ErrorCircuitSubmissionLimitReached)
		em.StatusCode = http.StatusBadRequest
	case ErrorRuleForOwnerNotFound:
		em.Msg = "Owner does not have associated rule."
		em.ErrCode = int(ErrorRuleForOwnerNotFound)
		em.StatusCode = http.StatusBadRequest
	case ErrorRuleNotFound:
		em.Msg = "Rule does not exist."
		em.ErrCode = int(ErrorRuleNotFound)
		em.StatusCode = http.StatusBadRequest
	case ErrorInvalidScore:
		em.Msg = "Failed to get simulation scores."
		em.ErrCode = int(ErrorInvalidScore)
		em.StatusCode = http.StatusInternalServerError
	case ErrorInvalidSummary:
		em.Msg = "Failed to get simulation summary."
		em.ErrCode = int(ErrorInvalidSummary)
		em.StatusCode = http.StatusInternalServerError
	case ErrorFailedToUploadLogs:
		em.Msg = "Failed to upload simulation logs."
		em.ErrCode = int(ErrorFailedToUploadLogs)
		em.StatusCode = http.StatusInternalServerError
	case ErrorCreditsExceeded:
		em.Msg = "Circuit credit limit exceeded."
		em.ErrCode = int(ErrorCreditsExceeded)
		em.StatusCode = http.StatusBadRequest
	case ErrorCircuitRuleNotFound:
		em.Msg = "Circuit rule not found."
		em.ErrCode = int(ErrorCircuitRuleNotFound)
		em.StatusCode = http.StatusBadRequest
	case ErrorFailedToGetLiveLogs:
		em.Msg = "Failed to get simulation live logs."
		em.ErrCode = int(ErrorFailedToGetLiveLogs)
		em.StatusCode = http.StatusInternalServerError
	case ErrorRobotIdentifierNotFound:
		em.Msg = "Robot identifier not found."
		em.ErrCode = int(ErrorRobotIdentifierNotFound)
		em.StatusCode = http.StatusNotFound
	case ErrorCompetitionNotStarted:
		em.Msg = "Competition has not started yet."
		em.ErrCode = int(ErrorCompetitionNotStarted)
		em.StatusCode = http.StatusInternalServerError
	case ErrorNotQualified:
		em.Msg = "Not qualified to compete on this circuit."
		em.ErrCode = int(ErrorNotQualified)
		em.StatusCode = http.StatusUnauthorized
	case ErrorLaunchHeldSimulation:
		em.Msg = "Failed to launch a held simulation."
		em.ErrCode = int(ErrorLaunchHeldSimulation)
		em.StatusCode = http.StatusInternalServerError
	case ErrorInvalidRobotImage:
		em.Msg = "Attempted to use a robot image that does not belong to the owner."
		em.ErrCode = int(ErrorInvalidRobotImage)
		em.StatusCode = http.StatusBadRequest
	case ErrorInvalidMarsupialSpecification:
		em.Msg = "Invalid marsupial specification. A parent and child must be specified, separated by a colon."
		em.ErrCode = int(ErrorInvalidMarsupialSpecification)
		em.StatusCode = http.StatusBadRequest
	case ErrorRobotModelLimitReached:
		em.Msg = "Too many robots of single model."
		em.ErrCode = int(ErrorRobotModelLimitReached)
		em.StatusCode = http.StatusBadRequest
	case ErrorSubmissionDeadlineReached:
		em.Msg = "Submission deadline reached."
		em.ErrCode = int(ErrorSubmissionDeadlineReached)
		em.StatusCode = http.StatusBadRequest
	}

	em.BaseError = errors.New(em.Msg)
	return em
}
