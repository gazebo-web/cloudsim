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
	// HasIPStatusCondition is used to indicate that pods have ips available.
	HasIPStatusCondition = Condition{
		Type:   "HasIPStatus",
		Status: "True",
	}
)
