package network

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNetworkPolicies(t *testing.T) {
	suite.Run(t, new(networkPoliciesTestSuite))
}

type networkPoliciesTestSuite struct {
	suite.Suite
	pod             *apiv1.Pod
	client          *fake.Clientset
	logger          ign.Logger
	networkPolicies *networkPolicies
}

func (s *networkPoliciesTestSuite) SetupTest() {
	s.pod = &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
			},
		},
		Spec:   apiv1.PodSpec{},
		Status: apiv1.PodStatus{},
	}
	s.client = fake.NewSimpleClientset()
	s.logger = ign.NewLoggerNoRollbar("TestNetworkPolicies", ign.VerbosityDebug)
	s.networkPolicies = &networkPolicies{
		API:    s.client,
		Logger: s.logger,
	}
}

func (s *networkPoliciesTestSuite) TestCreateEgressSpec() {
	egressRule := orchestrator.NetworkEgressRule{
		Ports:         []int32{1111, 2222, 3333},
		IPBlocks:      []string{"10.0.0.3/24"},
		AllowOutbound: true,
	}
	labels := map[string]string{
		"app2": "test",
	}
	to := orchestrator.NewSelector(labels)
	output := s.networkPolicies.createEgressSpec(egressRule, []orchestrator.Selector{to})
	s.Len(output, 6)
	for i, r := range output {
		switch i {
		case 0:
			s.Equal(int32(1111), r.Ports[0].Port.IntVal)
			break
		case 1:
			s.Equal(int32(2222), r.Ports[0].Port.IntVal)
			break
		case 2:
			s.Equal(int32(3333), r.Ports[0].Port.IntVal)
			break
		case 3:
			s.Equal("10.0.0.3/24", r.To[0].IPBlock.CIDR)
			break
		case 4:
			s.Equal(labels, r.To[0].PodSelector.MatchLabels)
			break
		case 5:
			s.NotNil(r)
			break
		}

	}
}

func (s *networkPoliciesTestSuite) TestCreateIngressSpec() {
	ingressRule := orchestrator.NetworkIngressRule{
		Ports:    []int32{1111, 2222, 3333},
		IPBlocks: []string{"10.0.0.3/24"},
	}
	labels := map[string]string{
		"app2": "test",
	}
	to := orchestrator.NewSelector(labels)
	output := s.networkPolicies.createIngressSpec(ingressRule, []orchestrator.Selector{to})
	s.Len(output, 5)
	for i, r := range output {
		switch i {
		case 0:
			s.Equal(int32(1111), r.Ports[0].Port.IntVal)
			break
		case 1:
			s.Equal(int32(2222), r.Ports[0].Port.IntVal)
			break
		case 2:
			s.Equal(int32(3333), r.Ports[0].Port.IntVal)
			break
		case 3:
			s.Equal("10.0.0.3/24", r.From[0].IPBlock.CIDR)
			break
		case 4:
			s.Equal(labels, r.From[0].PodSelector.MatchLabels)
			break
		}
	}
}

func (s *networkPoliciesTestSuite) TestCreateNetworkPolicy() {
	res, err := s.networkPolicies.Create(orchestrator.CreateNetworkPolicyInput{
		Name:      "test-np",
		Namespace: "default",
		Labels: map[string]string{
			"app": "test",
			"np":  "true",
		},
		PodSelector: orchestrator.NewSelector(s.pod.Labels),
		PeersFrom: []orchestrator.Selector{
			orchestrator.NewSelector(map[string]string{
				"app": "test",
			}),
		},
		PeersTo: []orchestrator.Selector{
			orchestrator.NewSelector(map[string]string{
				"app": "test",
			}),
		},
		Ingresses: orchestrator.NetworkIngressRule{
			Ports:    []int32{1111, 2222, 3333},
			IPBlocks: []string{"10.0.0.3/24"},
		},
		Egresses: orchestrator.NetworkEgressRule{
			Ports:         []int32{1111, 2222, 3333},
			IPBlocks:      []string{"10.0.0.3/24"},
			AllowOutbound: true,
		},
	})
	s.NoError(err)
	s.Equal("test-np", res.Name())
	s.Equal("default", res.Namespace())
	s.Equal(map[string]string{
		"app": "test",
		"np":  "true",
	}, res.Selector().Map())

	np, err := s.client.NetworkingV1().NetworkPolicies(res.Namespace()).Get("test-np", metav1.GetOptions{})
	s.NoError(err)
	s.Equal(res.Name(), np.Name)
}

func (s *networkPoliciesTestSuite) TestRemoveNetworkPolicy() {
	// Create a network policy
	res, err := s.networkPolicies.Create(orchestrator.CreateNetworkPolicyInput{
		Name:      "test-np",
		Namespace: "default",
		Labels: map[string]string{
			"app": "test",
			"np":  "true",
		},
		PodSelector: orchestrator.NewSelector(s.pod.Labels),
		PeersFrom: []orchestrator.Selector{
			orchestrator.NewSelector(map[string]string{
				"app": "test",
			}),
		},
		PeersTo: []orchestrator.Selector{
			orchestrator.NewSelector(map[string]string{
				"app": "test",
			}),
		},
		Ingresses: orchestrator.NetworkIngressRule{
			Ports:    []int32{1111, 2222, 3333},
			IPBlocks: []string{"10.0.0.3/24"},
		},
		Egresses: orchestrator.NetworkEgressRule{
			Ports:         []int32{1111, 2222, 3333},
			IPBlocks:      []string{"10.0.0.3/24"},
			AllowOutbound: true,
		},
	})

	s.Require().NoError(err)

	err = s.networkPolicies.Remove(res.Name(), res.Namespace())
	s.Assert().NoError(err)

	err = s.networkPolicies.Remove(res.Name(), res.Namespace())
	s.Assert().Error(err)
}
