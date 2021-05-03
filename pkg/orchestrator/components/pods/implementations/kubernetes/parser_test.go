package kubernetes

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource/phase"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	suite.Run(t, &ParserTestSuite{})
}

type ParserTestSuite struct {
	suite.Suite
}

func (s *ParserTestSuite) getBaseKubernetesPod(name, namespace string, phase phase.Phase) corev1.Pod {
	deletionTimestamp := metav1.NewTime(time.Now())
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: metav1.NewTime(time.Now()),
			DeletionTimestamp: &deletionTimestamp,
			Labels: map[string]string{
				"test": "true",
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPhase(phase),
		},
	}
}

func (s *ParserTestSuite) TestkubernetesPodToPodResource() {
	name := "test-name"
	namespace := "test-namespace"
	phase := phase.Running

	// Prepare a Kubernetes pod
	pod := s.getBaseKubernetesPod(name, namespace, phase)

	// Convert the Kubernetes pod to PodResource
	podRes := kubernetesPodToPodResource(pod)

	// Validate Resource
	s.Require().Equal(name, podRes.Name())
	s.Require().Equal(namespace, podRes.Namespace())
	s.Require().NotEmpty(podRes.Selector())
	// Validate ResourcePhase
	s.Require().Equal(phase, podRes.Phase())
	// Validate ResourceTimestamp
	s.Require().NotZero(podRes.CreationTimestamp())
	s.Require().NotZero(podRes.DeletionTimestamp())
}
