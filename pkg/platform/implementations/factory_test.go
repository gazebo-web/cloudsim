package implementations

import (
	"github.com/stretchr/testify/suite"
	email "gitlab.com/ignitionrobotics/web/cloudsim/pkg/email/implementations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	machines "gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines/implementations"
	configurationsImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations/implementations"
	ingressesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations"
	nodesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations"
	podsImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations"
	servicesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations"
	orchestrator "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	secrets "gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets/implementations"
	storage "gitlab.com/ignitionrobotics/web/cloudsim/pkg/storage/implementations"
	store "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	dependencies := factory.Dependencies{
		"Logger": logger,
	}

	var out platform.Platform
	s.Require().Nil(Factory.New(&config, dependencies, &out))
	s.Require().NotNil(out)
}
