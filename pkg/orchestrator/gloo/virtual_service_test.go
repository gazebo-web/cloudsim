package gloo

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestNewVirtualService(t *testing.T) {
	var vss orchestrator.Ingresses

	var gw gateway.GatewayV1Interface
	var logger ign.Logger

	vss = NewVirtualServices(gw, logger)

	assert.NotNil(t, vss)
	assert.IsType(t, &virtualServices{}, vss)
}
