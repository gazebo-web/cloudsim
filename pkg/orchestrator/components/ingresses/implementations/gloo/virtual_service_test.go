package gloo

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestNewVirtualService(t *testing.T) {
	var vss ingresses.Ingresses

	var gw gateway.GatewayV1Interface
	var logger ign.Logger
	var client *gloo.GlooV1Client

	vss = NewVirtualServices(gw, logger, client)

	assert.NotNil(t, vss)
	assert.IsType(t, &VirtualServices{}, vss)
}
