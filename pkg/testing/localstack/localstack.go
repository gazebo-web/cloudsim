package localstack

type LocalStack interface {
}

type localStack struct {
	host     string
	edgePort int
}

func New() LocalStack {
	return &localStack{}
}
