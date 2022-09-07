package resource

// Condition represents a state that should be reached.
type Condition struct {
	Type   string
	Status string
}

var (
	// ReadyCondition is used to indicate that Nodes and Pods are ready.
	ReadyCondition = Condition{
		Type:   "Ready",
		Status: "True",
	}
	// SucceededCondition is used to indicate that Pods have completed successfully and returned no errors.
	SucceededCondition = Condition{
		Type:   "Succeeded",
		Status: "True",
	}
	// FailedCondition is used to indicate that Pods have failed.
	FailedCondition = Condition{
		Type:   "Failed",
		Status: "True",
	}
	// HasIPStatusCondition is used to indicate that pods have ips available.
	HasIPStatusCondition = Condition{
		Type:   "HasIPStatus",
		Status: "True",
	}
)
