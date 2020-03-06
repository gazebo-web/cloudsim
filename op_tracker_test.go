package main

// OpType represents the identifier of an operation to invoke.
type OpType string

// OpTracker is an object composed by mocks that allow operation tracking.
type OpTracker struct {
	invokedOps        []OpType
	invokedOpsResults []interface{}
}

// NewOpTracker creates a new OpTracker.
func NewOpTracker() *OpTracker {
	tracker := &OpTracker{}
	tracker.Reset()
	return tracker
}

// Reset sets all values to their initial configuration.
func (m *OpTracker) Reset() {
	m.invokedOps = make([]OpType, 0)
	m.invokedOpsResults = make([]interface{}, 0)
}

// CountCalls returns how many times a given operation was called since the last Reset.
func (m *OpTracker) CountCalls(op OpType) int {
	c := 0
	for _, inv := range m.invokedOps {
		if inv == op {
			c++
		}
	}
	return c
}

// TrackCall registers that an operation was called.
func (m *OpTracker) TrackCall(op OpType) {
	m.invokedOps = append(m.invokedOps, op)
}

// TrackCallResult stores the value returned by an operation call.
func (m *OpTracker) TrackCallResult(res interface{}) {
	m.invokedOpsResults = append(m.invokedOpsResults, res)
}
