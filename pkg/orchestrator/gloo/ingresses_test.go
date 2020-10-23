package gloo

import (
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"testing"
)

func TestNewIngresses(t *testing.T) {
	var api gloo.GlooV1Interface
	var gw gateway.GatewayV1Interface
	var ingress orchestrator.Ingresses

	ingress = NewIngresses(api, gw)

	assert.NotNil(t, ingress)
}
