package jobs

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestCreateConfigurationsSuite(t *testing.T) {
	suite.Run(t, new(testCreateConfigurationsSuite))
}

type testCreateConfigurationsSuite struct {
	suite.Suite
	kubernetesAPI      *fake.Clientset
	cluster            orchestrator.Cluster
	jobName            string
	namespace          string
	configurationName1 string
	configurationName2 string
	configurationName3 string
}

func (suite *testCreateConfigurationsSuite) SetupSuite() {
	logger := ign.NewLoggerNoRollbar("TestCreateConfigurationsSuite", ign.VerbosityInfo)
	suite.cluster, suite.kubernetesAPI = kubernetes.NewFakeKubernetes(logger)

	suite.jobName = "test_job"

	suite.namespace = "default"

	suite.configurationName1 = "test-1"
	suite.configurationName2 = "test-2"
	suite.configurationName3 = "test-3"

	_, err := suite.cluster.Configurations().Create(
		configurations.CreateConfigurationInput{
			Name:      suite.configurationName1,
			Namespace: suite.namespace,
		},
	)
	suite.Require().NoError(err)

	_, err = suite.cluster.Configurations().Create(
		configurations.CreateConfigurationInput{
			Name:      suite.configurationName2,
			Namespace: suite.namespace,
		},
	)
	suite.Require().NoError(err)

	_, err = suite.cluster.Configurations().Create(
		configurations.CreateConfigurationInput{
			Name:      suite.configurationName3,
			Namespace: suite.namespace,
		},
	)
	suite.Require().NoError(err)

}

func (suite *testCreateConfigurationsSuite) getNumberOfConfigurations() int {
	cms, err := suite.kubernetesAPI.CoreV1().ConfigMaps(suite.namespace).List(metav1.ListOptions{})
	suite.Require().NoError(err)

	return len(cms.Items)
}

func (suite *testCreateConfigurationsSuite) TestRemoveCreatedConfigurationsOnFailureRollbackHandler() {
	// Get DB
	db, err := gormUtils.GetTestDBFromEnvVars()
	suite.Require().NoError(err)

	err = actions.CleanAndMigrateDB(db)
	suite.Require().NoError(err)

	// Create action to register the job's datatypes in the registry
	_, err = actions.NewAction(
		actions.Jobs{
			CreateConfigurations,
		},
	)
	suite.Require().NoError(err)

	// Create store
	p, err := platform.NewPlatform(
		"test", platform.Components{
			Cluster: suite.cluster,
		},
	)
	suite.Require().NoError(err)
	state := &TestState{
		platform: p,
	}
	store := state.ToStore()

	// Create job data
	deployment := &actions.Deployment{
		Action:     suite.jobName,
		CurrentJob: suite.jobName,
	}
	suite.Require().NoError(db.Create(deployment).Error)

	err = deployment.SetJobData(
		db, &suite.jobName, actions.DeploymentJobInput, &CreateConfigurationsInput{
			{
				Namespace: suite.namespace,
				Name:      suite.configurationName1,
			},
			{
				Namespace: suite.namespace,
				Name:      suite.configurationName2,
			},
			{
				Namespace: suite.namespace,
				Name:      suite.configurationName3,
			},
		},
	)
	suite.Require().NoError(err)

	// Check that configurations exist
	suite.Require().Equal(3, suite.getNumberOfConfigurations())

	expectedErr := errors.New("error")
	_, err = DeleteCreatedConfigurationsOnFailure(store, db, deployment, nil, expectedErr)
	suite.Equal(expectedErr, err)

	// Verify that configurations no longer exist
	suite.Require().Equal(0, suite.getNumberOfConfigurations())
}
