package gloo

import (
	"errors"
	"fmt"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Gloo struct {
	client gloo.GlooV1Client
	logger ign.Logger
}

func (g *Gloo) GetUpstream(namespace string, selector orchestrator.Selector) (orchestrator.Resource, error) {
	g.logger.Debug(
		fmt.Sprintf("Getting upstream on namespace [%s] pointing to the given labels [%s]",
			namespace, selector.Map()),
	)

	list, err := g.client.Upstreams(namespace).List(v1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		g.logger.Debug(
			fmt.Sprintf("Failed to get upstream on namespace [%s] pointing to the given labels [%s]. Error: %s",
				namespace, selector.Map()),
		)
		return nil, err
	}

	if len(list.Items) < 1 {
		return nil, errors.New("did not find a Gloo upstream for target Kubernetes service")
	} else if len(list.Items) > 1 {
		return nil, errors.New("found too many Gloo upstreams for target Kubernetes service")
	}

	s := orchestrator.NewSelector(list.Items[0].Labels)
	res := orchestrator.NewResource(list.Items[0].Name, namespace, s)
	return res, nil
}

func NewExtensions(client gloo.GlooV1Client, logger ign.Logger) orchestrator.Extensions {
	return &Gloo{
		client: client,
		logger: logger,
	}
}
