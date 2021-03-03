package testing

import (
	"k8s.io/apimachinery/pkg/runtime"
	k8testing "k8s.io/client-go/testing"
)

// ChainedMockFunction is a generic function used to create Kubernetes Reactor
// chains. Reactor chains allows modifying resources in response to events.
// ChainedMockFunction differs from MockFunction in that it returns a `handled`
// value, which lets the reactor chain short circuit if necessary.
type ChainedMockFunction func(args ...interface{}) (handled bool, res interface{})

// GenerateReactor generates a reaction function that can be used in Kubernetes fake implementations.
// The fn function performs the actual reaction operation. This function can short circuit the reactor chain by way of
// the handled return value.
func GenerateReactor(fn ChainedMockFunction) k8testing.ReactionFunc {
	return func(action k8testing.Action) (handled bool, ret runtime.Object, err error) {
		handled, res := fn(action)
		// If the mock result is an error, return that error
		if err, ok := res.(error); ok {
			return handled, nil, err
		}
		ret = res.(runtime.Object)
		return handled, ret, nil
	}
}

// GenerateObjectReturnChainedMock is a helper function to generate a ChainedMockFunction that returns a provided
// object without performing any additional operations. Typically used to return a specific value for a set of
// Kubernetes fake verb and resource combination.
func GenerateObjectReturnChainedMock(handled bool, res interface{}) ChainedMockFunction {
	return func(args ...interface{}) (bool, interface{}) {
		return handled, res
	}
}
