package localstack

// LocalStack groups a set of methods to wrap LocalStack tool for testing purposes.
type LocalStack interface {
}

// localStack is a LocalStack implementation.
type localStack struct {
	host     string
	edgePort int
}

// New initializes a new LocalStack implementation.
func New() LocalStack {
	return &localStack{}
}
