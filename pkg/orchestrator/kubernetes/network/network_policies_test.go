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
