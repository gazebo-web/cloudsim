package gloo

import (
	"context"
	"fmt"
	gatewayapiv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGlooVirtualServiceTestSuite(t *testing.T) {
	suite.Run(t, &GlooVirtualServiceTestSuite{})
}

// GlooVirtualServiceTestSuite tests Gloo Virtual Service operations.
type GlooVirtualServiceTestSuite struct {
	suite.Suite
	ctx                    context.Context
	gloo                   Clientset
	virtualServiceManifest *gatewayv1.VirtualService
	virtualService         *gatewayv1.VirtualService
	virtualServiceRoutes   []*gatewayapiv1.Route
	virtualHostDomain      string
	virtualHostRouteName   string
}

func (suite *GlooVirtualServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.gloo = newFakeClientset()
	suite.virtualHostDomain = "test-service.net"

	// The virtualService contains four dummy routes and a target route
	// The dummy routes are set in place to make sure that operations only modify target routes
	suite.virtualServiceRoutes = []*gatewayapiv1.Route{
		suite.generateDummyRoute(1),
		suite.generateTestVirtualHostRoute(1),
		suite.generateDummyRoute(2),
		suite.generateTestVirtualHostRoute(1),
		suite.generateDummyRoute(3),
	}
	suite.virtualServiceManifest = &gatewayv1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: gatewayapiv1.VirtualService{
			DisplayName: "test",
			VirtualHost: &gatewayapiv1.VirtualHost{
				Domains: []string{suite.virtualHostDomain},
				Routes:  suite.virtualServiceRoutes,
			},
		},
	}
}

func (suite *GlooVirtualServiceTestSuite) generateTestVirtualHostRoute(i int) *gatewayapiv1.Route {
	return CreateVirtualHostRoute(
		suite.virtualHostRouteName,
		[]*matchers.Matcher{
			CreateVirtualHostRouteExactMatcher(fmt.Sprintf("/test/%d", i)),
		},
		CreateVirtualHostRouteAction(metav1.NamespaceDefault, "test"),
	)
}

func (suite *GlooVirtualServiceTestSuite) generateDummyRoute(i int) *gatewayapiv1.Route {
	return CreateVirtualHostRoute(
		fmt.Sprintf("dummy-%d", i),
		[]*matchers.Matcher{
			CreateVirtualHostRouteExactMatcher(fmt.Sprintf("/dummy/%d", i)),
		},
		CreateVirtualHostRouteAction(metav1.NamespaceDefault, "test"),
	)
}

func (suite *GlooVirtualServiceTestSuite) SetupTest() {
	var err error
	suite.virtualService, err = suite.gloo.Gateway().VirtualServices(metav1.NamespaceDefault).Create(
		suite.virtualServiceManifest,
	)
	suite.NoError(err)
}

func (suite *GlooVirtualServiceTestSuite) TearDownTest() {
	err := suite.gloo.Gateway().VirtualServices(metav1.NamespaceDefault).Delete(
		suite.virtualServiceManifest.Name, nil,
	)
	suite.NoError(err)
}

func (suite *GlooVirtualServiceTestSuite) TestGetVirtualService() {
	virtualService, err := getVirtualService(
		suite.ctx,
		suite.gloo,
		metav1.NamespaceDefault,
		suite.virtualService.Name,
	)
	suite.NoError(err)
	suite.Equal(suite.virtualService, virtualService)
}

func (suite *GlooVirtualServiceTestSuite) TestUpsertVirtualServiceRouteRemoveVirtualServiceRoute() {
	// Check that routes in a specific VirtualService index match the specified singleRoute
	test := func(virtualService *gatewayv1.VirtualService, routes ...*gatewayapiv1.Route) {
		for i, route := range routes {
			vhRoute := virtualService.Spec.VirtualHost.Routes[i]
			suite.Equal(vhRoute, route)
		}
	}

	upsert := func(virtualService *gatewayv1.VirtualService, routes ...*gatewayapiv1.Route) *gatewayv1.VirtualService {
		virtualService, err := UpsertVirtualServiceRoute(
			suite.ctx,
			suite.gloo,
			metav1.NamespaceDefault,
			virtualService.Name,
			routes...,
		)
		suite.NoError(err)
		return virtualService
	}

	remove := func(virtualService *gatewayv1.VirtualService, routes ...*gatewayapiv1.Route) *gatewayv1.VirtualService {
		virtualService, err := RemoveVirtualServiceRoute(
			suite.ctx,
			suite.gloo,
			metav1.NamespaceDefault,
			virtualService.Name,
			routes...,
		)
		suite.NoError(err)
		return virtualService
	}

	// Update existent route
	existingRoute := CreateVirtualHostRoute(
		suite.virtualServiceRoutes[1].Name,
		[]*matchers.Matcher{},
		CreateVirtualHostRouteAction(metav1.NamespaceDefault, "single"),
	)
	virtualService := upsert(suite.virtualService, existingRoute)
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
	)

	// Add a single route
	singleRoute := CreateVirtualHostRoute(
		"single",
		[]*matchers.Matcher{},
		CreateVirtualHostRouteAction(metav1.NamespaceDefault, "single"),
	)
	virtualService = upsert(suite.virtualService, singleRoute)
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
		singleRoute,
	)

	// Add multiple routes
	multiRoutes := []*gatewayapiv1.Route{
		CreateVirtualHostRoute(
			"multi-1",
			[]*matchers.Matcher{},
			CreateVirtualHostRouteAction(metav1.NamespaceDefault, "multi-1"),
		),
		CreateVirtualHostRoute(
			"multi-2",
			[]*matchers.Matcher{},
			CreateVirtualHostRouteAction(metav1.NamespaceDefault, "multi-2"),
		),
		CreateVirtualHostRoute(
			"multi-3",
			[]*matchers.Matcher{},
			CreateVirtualHostRouteAction(metav1.NamespaceDefault, "multi-3"),
		),
	}
	virtualService = upsert(suite.virtualService, multiRoutes...)
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
		singleRoute,
		multiRoutes[0],
		multiRoutes[1],
		multiRoutes[2],
	)

	// Remove a non-existent singleRoute, and check that nothing changed
	nonexistentRoute := CreateVirtualHostRoute("nonexistent", nil, nil)
	virtualService = remove(virtualService, nonexistentRoute)
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
		singleRoute,
		multiRoutes[0],
		multiRoutes[1],
		multiRoutes[2],
	)

	// Remove a single non-dummy singleRoute
	virtualService = remove(virtualService, multiRoutes[0])
	// The deletion process swaps deleted entries with the last entry in the array
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
		singleRoute,
		multiRoutes[2],
		multiRoutes[1],
	)

	// Remove multiple non-dummy routes
	virtualService = remove(virtualService, singleRoute, nonexistentRoute, multiRoutes[1])
	test(
		virtualService,
		suite.virtualServiceRoutes[0],
		existingRoute,
		suite.virtualServiceRoutes[2],
		suite.virtualServiceRoutes[3],
		suite.virtualServiceRoutes[4],
		multiRoutes[2],
	)
}
