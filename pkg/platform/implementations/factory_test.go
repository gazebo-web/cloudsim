package implementations

import (
	email "github.com/gazebo-web/cloudsim/pkg/email/implementations"
	"github.com/gazebo-web/cloudsim/pkg/factory"
	machines "github.com/gazebo-web/cloudsim/pkg/machines/implementations"
	configurationsImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations/implementations"
	ingressesImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/network/implementations"
	nodesImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/nodes/implementations"
	podsImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/pods/implementations"
	servicesImpl "github.com/gazebo-web/cloudsim/pkg/orchestrator/components/services/implementations"
	orchestrator "github.com/gazebo-web/cloudsim/pkg/orchestrator/implementations"
	"github.com/gazebo-web/cloudsim/pkg/platform"
	secrets "github.com/gazebo-web/cloudsim/pkg/secrets/implementations"
	storage "github.com/gazebo-web/cloudsim/pkg/storage/implementations"
	store "github.com/gazebo-web/cloudsim/pkg/store/implementations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestPlatformFactorySuite(t *testing.T) {
	suite.Run(t, new(testPlatformFactorySuite))
}

type testPlatformFactorySuite struct {
	suite.Suite
}

func (s *testPlatformFactorySuite) TestNewFunc() {
	// Prepare config
	// Use default kubeconfig
	kubeconfig := factory.ConfigValues{
		"kubeconfig": "",
	}
	// Default to us-east-1
	defaultRegion := "us-east-1"
	config := factory.Config{
		Type: Default,
		Config: factory.ConfigValues{
			"Name": "us-east-1",
			"Components": factory.ConfigValues{
				"Machines": factory.ConfigValues{
					"type": machines.EC2,
					"config": factory.ConfigValues{
						"region": defaultRegion,
						"zones": []factory.ConfigValues{
							{
								"zone":     "us-east-1a",
								"subnetID": "sg-123456789",
							},
						},
					},
				},
				"Orchestrator": factory.ConfigValues{
					"type": orchestrator.Kubernetes,
					"config": factory.ConfigValues{
						"api": kubeconfig,
						"components": factory.ConfigValues{
							"nodes":           factory.ConfigValues{"type": nodesImpl.Kubernetes},
							"pods":            factory.ConfigValues{"type": podsImpl.Kubernetes},
							"services":        factory.ConfigValues{"type": servicesImpl.Kubernetes},
							"ingresses":       factory.ConfigValues{"type": ingressesImpl.Kubernetes},
							"ingressRules":    factory.ConfigValues{"type": ingressesImpl.Kubernetes},
							"networkPolicies": factory.ConfigValues{"type": networkImpl.Kubernetes},
							"configurations":  factory.ConfigValues{"type": configurationsImpl.Kubernetes},
						},
					},
				},
				"Storage": factory.ConfigValues{
					"type": storage.S3,
					"config": factory.ConfigValues{
						"region": defaultRegion,
					},
				},
				"Store": factory.ConfigValues{
					"type": store.Store,
					"config": factory.ConfigValues{
						"machinesStore": factory.ConfigValues{
							"instanceProfileValue": "test",
							"keyNameValue":         "test",
							"namePrefixValue":      "test",
							"clusterNameValue":     "test",
						},
						"ignitionStore": factory.ConfigValues{
							"defaultSenderValue": "test@ignitionrobotics.org",
							"logsBucketValue":    "test_bucket",
						},
						"orchestratorStore": factory.ConfigValues{
							"ingressNameValue": "test",
							"ingressHostValue": "test.com",
						},
					},
				},
				"Secrets": factory.ConfigValues{
					"type": secrets.Kubernetes,
					"config": factory.ConfigValues{
						"api": kubeconfig,
					},
				},
				"EmailSender": factory.ConfigValues{
					"type": email.SES,
					"config": factory.ConfigValues{
						"region": defaultRegion,
					},
				},
			},
		},
	}

	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	dependencies := factory.Dependencies{
		"Logger": logger,
	}

	var out platform.Platform
	s.Require().Nil(Factory.New(&config, dependencies, &out))
	s.Require().NotNil(out)
}
