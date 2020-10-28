package gloo

import (
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestNewVirtualHosts(t *testing.T) {
	var vhs orchestrator.IngressRules
	var gw gateway.GatewayV1Interface
	var logger ign.Logger
	vhs = NewVirtualHosts(gw, logger)
	assert.IsType(t, &virtualHosts{}, vhs)
}

func newTestVirtualService() *gatewayv1.VirtualService {
	return &gatewayv1.VirtualService{Spec: v1.VirtualService{
		VirtualHost: &v1.VirtualHost{
			Domains: []string{"openrobotics.org"},
			Routes: []*v1.Route{
				{
					Matchers:             []*matchers.Matcher{},
					Action:               nil,
					Options:              nil,
					Name:                 "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			Options:              nil,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		SslConfig:            nil,
		DisplayName:          "test",
		Status:               core.Status{},
		Metadata:             core.Metadata{},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}}
}