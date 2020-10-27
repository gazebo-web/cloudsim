package gloo

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
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
