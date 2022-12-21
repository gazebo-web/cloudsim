package kubernetes

import (
	"context"
	"fmt"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestConfigMaps(t *testing.T) {
	suite.Run(t, new(configMapsTestSuite))
}

type configMapsTestSuite struct {
	suite.Suite
	pod        *apiv1.Pod
	client     *fake.Clientset
	logger     gz.Logger
	configMaps *configMaps
}

func (s *configMapsTestSuite) SetupTest() {
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
	s.logger = gz.NewLoggerNoRollbar("TestConfigMaps", gz.VerbosityDebug)
	s.configMaps = &configMaps{
		API:    s.client,
		Logger: s.logger,
	}
}

func (s *configMapsTestSuite) TestCreateConfigMap() {
	res, err := s.configMaps.Create(context.TODO(), configurations.CreateConfigurationInput{
		Name:      "test-np",
		Namespace: "default",
		Labels: map[string]string{
			"app": "test",
			"np":  "true",
		},
	})
	s.NoError(err)
	s.Equal("test-np", res.Name())
	s.Equal("default", res.Namespace())
	s.Equal(
		map[string]string{
			"app": "test",
			"np":  "true",
		}, res.Selector().Map(),
	)

	np, err := s.client.CoreV1().ConfigMaps(res.Namespace()).Get(context.TODO(), "test-np", metav1.GetOptions{})
	s.NoError(err)
	s.Equal(res.Name(), np.Name)
}

func (s *configMapsTestSuite) TestDeleteConfiguration() {
	// Create a config map
	res, err := s.configMaps.Create(context.TODO(), configurations.CreateConfigurationInput{
		Name:      "test-np",
		Namespace: "default",
		Labels: map[string]string{
			"app": "test",
			"np":  "true",
		},
	})

	s.Require().NoError(err)

	// Trying to remove an existent config map should not fail
	_, err = s.configMaps.Delete(context.TODO(), res)
	s.Assert().NoError(err)

	// Trying to remove a nonexistent config map should fail
	_, err = s.configMaps.Delete(context.TODO(), res)
	s.Assert().Error(err)
}

func (s *configMapsTestSuite) setupTestRemoveBulk() {
	for i := 0; i < 3; i++ {
		_, err := s.configMaps.Create(context.TODO(), configurations.CreateConfigurationInput{
			Name:      fmt.Sprintf("test-np-%d", i),
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
				"np":  fmt.Sprintf("%d", i),
			},
		})
		s.Require().NoError(err)
	}
}
