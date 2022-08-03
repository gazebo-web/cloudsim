package jobs

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	fakeStore "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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
	store              store.Store
	jobName            string
	namespace          string
	configurationName1 string
	configurationName2 string
	configurationName3 string
}

func (suite *testCreateConfigurationsSuite) SetupSuite() {
	logger := ign.NewLoggerNoRollbar("TestCreateConfigurationsSuite", ign.VerbosityInfo)
	suite.cluster, suite.kubernetesAPI = kubernetes.NewFakeKubernetes(logger)

	suite.store = fakeStore.NewDefaultFakeStore()
	// Set up store
	storeIgnition := fakeStore.NewFakeIgnition()
	storeOrchestrator := fakeStore.NewFakeOrchestrator()

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")

	suite.store = fakeStore.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	suite.namespace = "default"

	suite.configurationName1 = "test-1"
	suite.configurationName2 = "test-2"
	suite.configurationName3 = "test-3"
}

func (suite *testCreateConfigurationsSuite) getNumberOfConfigurations() int {
	cms, err := suite.kubernetesAPI.CoreV1().ConfigMaps(suite.namespace).List(context.TODO(), metav1.ListOptions{})
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
			Store:   suite.store,
		},
	)
	suite.Require().NoError(err)
	state := &TestState{
		platform: p,
	}
	store := state.ToStore()

	// Create deployment
	deployment := &actions.Deployment{
		Action:     suite.jobName,
		CurrentJob: suite.jobName,
	}

	// There should be no pre-existing configurations
	suite.Require().Equal(0, suite.getNumberOfConfigurations())

	// Create the configurations
	_, err = CreateConfigurations.Run(store, db, deployment, CreateConfigurationsInput{
		{
			Name:      suite.configurationName1,
			Namespace: suite.namespace,
		},
		{
			Name:      suite.configurationName2,
			Namespace: suite.namespace,
		},
		{
			Name:      suite.configurationName3,
			Namespace: suite.namespace,
		},
	})
	suite.Require().NoError(err)

	// There should be 3 configurations
	suite.Require().Equal(3, suite.getNumberOfConfigurations())

	// Run the rollback handler
	err = errors.New("error")
	_, err = DeleteCreatedConfigurationsOnFailure(store, db, deployment, nil, err)
	suite.Assert().NoError(err)

	// Verify that configurations no longer exist
	suite.Require().Equal(0, suite.getNumberOfConfigurations())
}
