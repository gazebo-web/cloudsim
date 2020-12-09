package gloo

import (
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gatewayFake "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1/fake"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewVirtualHosts(t *testing.T) {
	var vhs orchestrator.IngressRules
	var gw gateway.GatewayV1Interface
	var logger ign.Logger
	vhs = NewVirtualHosts(gw, logger)
	assert.IsType(t, &virtualHosts{}, vhs)
}

func setupTestVirtualHosts(t *testing.T, regex string) (*gatewayFake.FakeGatewayV1, ign.Logger, orchestrator.Resource) {
	// Define constants
	const namespace = "default"
	const upstream = "my-service"

	domains := []string{"test.org"}

	// Initialize new mock virtual service
	vs := newTestVirtualService(t.Name(), namespace, upstream, regex, domains)

	// Initialize fake gateway implementation.
	gw := &gatewayFake.FakeGatewayV1{Fake: &fake.NewSimpleClientset().Fake}

	// Create the mock virtual service
	vs, err := gw.VirtualServices(namespace).Create(vs)
	require.NoError(t, err)

	// Initialize logger
	logger := ign.NewLoggerNoRollbar("TestVirtualHosts", ign.VerbosityDebug)

	// Initialize virtual services manager
	vss := NewVirtualServices(gw, logger)

	// Get the resource associated with the virtual service.
	res, err := vss.Get(t.Name(), namespace)
	require.NotNil(t, res)
	require.NoError(t, err)

	return gw, logger, res
}

func TestVirtualHosts_Get(t *testing.T) {
	const regex = "[a-zA-Z]+"
	gw, logger, res := setupTestVirtualHosts(t, regex)

	// Initialize virtual hosts manager
	vhs := NewVirtualHosts(gw, logger)

	// Get the representation of the virtual host that has the test.org domain.
	rule, err := vhs.Get(res, "test.org")
	require.NoError(t, err)

	// Get the paths from the rule
	p := rule.Paths()

	// There should only be one path, and it should have the values from the virtual host.
	assert.Len(t, p, 1)
	assert.Equal(t, t.Name(), p[0].UID)
	assert.Equal(t, regex, p[0].Address)

	// --------------------------------------------------------------------------

	// If we try to get an unknown domain, it will throw an error.
	_, err = vhs.Get(res, "another-test.org")
	require.Error(t, err)
}

func TestVirtualHosts_Upsert(t *testing.T) {
	const regex = "[a-zA-Z]+"
	gw, logger, res := setupTestVirtualHosts(t, regex)

	// Initialize virtual hosts manager
	vhs := NewVirtualHosts(gw, logger)

	// Get the representation of the virtual host that has the test.org domain.
	rule, err := vhs.Get(res, "test.org")
	require.NoError(t, err)

	p := NewPath(t.Name(), GenerateRegexMatcher("another-regex"), GenerateRouteAction("default", "my-new-service"))
	err = vhs.Upsert(rule, p)
	assert.NoError(t, err)

	rule, err = vhs.Get(res, "test.org")
	require.NoError(t, err)

	assert.Len(t, rule.Paths(), 1)

	assert.Equal(t, "another-regex", rule.Paths()[0].Address)
	assert.Equal(t, "my-new-service", rule.Paths()[0].Endpoint.Name)
}

func TestVirtualHost_Remove(t *testing.T) {
	const regex = "[a-zA-Z]+"
	gw, logger, res := setupTestVirtualHosts(t, regex)

	// Initialize virtual hosts manager
	vhs := NewVirtualHosts(gw, logger)

	// Get the representation of the virtual host that has the test.org domain.
	rule, err := vhs.Get(res, "test.org")
	require.NoError(t, err)

	p := rule.Paths()

	err = vhs.Remove(rule, p...)
	assert.NoError(t, err)

	rule, err = vhs.Get(res, "test.org")
	require.NoError(t, err)

	assert.Len(t, rule.Paths(), 0)

	err = vhs.Remove(rule, p...)
	assert.Error(t, err)
	assert.Equal(t, err, orchestrator.ErrRuleEmpty)
}

func newTestVirtualService(name, namespace, upstream, regex string, domains []string) *gatewayv1.VirtualService {
	return &gatewayv1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.VirtualService{
			VirtualHost: &v1.VirtualHost{
				Domains: domains,
				Routes: []*v1.Route{
					{
						Matchers: []*matchers.Matcher{
							GenerateRegexMatcher(regex),
						},
						Action: GenerateRouteAction(namespace, upstream),
						Name:   name,
					},
				},
			},
			DisplayName: name,
		},
	}
}
