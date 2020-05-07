package simulations

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"log"
)

// Deprecated: K8Factory is the single place where K8 instances are created.
type K8Factory struct {
	isGoTest               bool
	connectToCloudServices bool
}

// Deprecated: NewK8Factory creates a new K8 factory
func NewK8Factory(isGoTest, connectToCloudServices bool) *K8Factory {
	f := K8Factory{}
	f.isGoTest = isGoTest
	f.connectToCloudServices = connectToCloudServices
	return &f
}

// Deprecated: NewK8 creates a new instance of Kubernetes client.
func (f *K8Factory) NewK8(ctx context.Context) kubernetes.Interface {

	if f.isGoTest {
		return &MockableClientset{Interface: fake.NewSimpleClientset()}
	} else if f.connectToCloudServices {
		kcli, err := GetKubernetesClient()
		if err != nil {
			logger(ctx).Critical("Critical error trying to create a client to kubernetes", err)
			log.Fatalf("%+v\n", err)
		}
		return kcli
	}
	return nil
}

// Deprecated: MockableClientset is a type used in tests to allow for easy mocking of
// Kubernetes clientset.
type MockableClientset struct {
	kubernetes.Interface
}

// Deprecated: AssertMockedClientset casts the given arg to MockableClientset or fails.
func AssertMockedClientset(cli kubernetes.Interface) *MockableClientset {
	return cli.(*MockableClientset)
}

// SetImpl sets the underlying implementation of this MockableClientset
func (m *MockableClientset) SetImpl(cli kubernetes.Interface) {
	m.Interface = cli
}
