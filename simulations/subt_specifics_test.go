package simulations

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestIngressTestSuite(t *testing.T) {
	suite.Run(t, &IngressTestSuite{})
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type IngressTestSuite struct {
	suite.Suite
	ctx                       context.Context
	kcli                      *fake.Clientset
	ingressManifest           *v1beta1.Ingress
	ingress                   *v1beta1.Ingress
	ingressHost               string
	ingressAnonymousHost      string
	ingressRulePath           *v1beta1.IngressRuleValue
	ingressHostRuleIndex      int
	ingressAnonymousRuleIndex int
}

type IngressPath struct {
	path        string
	service     string
	servicePort int
}

func (suite *IngressTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.kcli = fake.NewSimpleClientset()
	suite.ingressHost = "test-service.net"
	suite.ingressAnonymousHost = ""
	// The ingress contains two dummy rules and a target rule
	// These dummy rules are there to check that they are not be modified
	suite.ingressManifest = &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test",
			Namespace: v1.NamespaceDefault,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				// 0: Dummy rule, this rule should never be updated.
				{
					Host: "dummy.net",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/dummy",
									Backend: v1beta1.IngressBackend{
										ServiceName: "dummy-1",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 1234,
										},
									},
								},
							},
						},
					},
				},
				// 1: Rule with host. This rule will be targeted by tests.
				{
					Host: suite.ingressHost,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/dummy",
									Backend: v1beta1.IngressBackend{
										ServiceName: "dummy-3",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3456,
										},
									},
								},
							},
						},
					},
				},
				// 2: Copy of the previous rule put in place to check that operations only affect the first rule
				// matching a host.
				{
					Host: suite.ingressHost,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/dummy",
									Backend: v1beta1.IngressBackend{
										ServiceName: "dummy-copy",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 9999,
										},
									},
								},
							},
						},
					},
				},
				// 3: Rule with anonymous host. This rule will be targeted by tests.
				{
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/dummy",
									Backend: v1beta1.IngressBackend{
										ServiceName: "dummy-2",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 2345,
										},
									},
								},
							},
						},
					},
				},
				// 4: Copy of the previous rule put in place to check that operations only affect the first rule
				// matching a host.
				{
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/dummy",
									Backend: v1beta1.IngressBackend{
										ServiceName: "dummy-copy",
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 7777,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	suite.ingressHostRuleIndex = 1
	suite.ingressAnonymousRuleIndex = 3
}

func (suite *IngressTestSuite) SetupTest() {
	var err error
	suite.ingress, err = suite.kcli.ExtensionsV1beta1().Ingresses(v1.NamespaceDefault).Create(suite.ingressManifest)
	suite.NoError(err)
}

func (suite *IngressTestSuite) TearDownTest() {
	err := suite.kcli.ExtensionsV1beta1().Ingresses(v1.NamespaceDefault).Delete(suite.ingressManifest.Name, nil)
	suite.NoError(err)
}

func (suite *IngressTestSuite) generateRulePath(path IngressPath) *v1beta1.HTTPIngressPath {
	return &v1beta1.HTTPIngressPath{
		Path: path.path,
		Backend: v1beta1.IngressBackend{
			ServiceName: path.service,
			ServicePort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(path.servicePort),
			},
		},
	}
}

func (suite *IngressTestSuite) TestGetIngress() {
	ingress, err := getIngress(
		suite.ctx,
		suite.kcli,
		v1.NamespaceDefault,
		suite.ingress.Name,
	)
	suite.NoError(err)
	suite.Equal(suite.ingress, ingress)
}

func (suite *IngressTestSuite) TestGetIngressRule() {
	// Get rule for specific host
	rule, err := getIngressRule(
		suite.ctx,
		suite.ingress,
		&suite.ingressHost,
	)
	suite.NoError(err)
	suite.Equal(suite.ingress.Spec.Rules[suite.ingressHostRuleIndex].HTTP, rule)

	// Get rule for anonymous host
	rule, err = getIngressRule(
		suite.ctx,
		suite.ingress,
		&suite.ingressAnonymousHost,
	)
	suite.NoError(err)
	suite.Equal(suite.ingress.Spec.Rules[suite.ingressAnonymousRuleIndex].HTTP, rule)
}

func (suite *IngressTestSuite) TestUpsertIngressRuleRemoveIngressRule() {
	upsert := func(ingress *v1beta1.Ingress, paths ...*v1beta1.HTTPIngressPath) *v1beta1.Ingress {
		ingress, err := upsertIngressRule(
			suite.ctx,
			suite.kcli,
			v1.NamespaceDefault,
			ingress.Name,
			&suite.ingressHost,
			paths...,
		)
		suite.NoError(err)
		return ingress
	}

	remove := func(ingress *v1beta1.Ingress, paths ...string) *v1beta1.Ingress {
		ingress, err := removeIngressRule(
			suite.ctx,
			suite.kcli,
			v1.NamespaceDefault,
			ingress.Name,
			&suite.ingressHost,
			paths...,
		)
		suite.NoError(err)
		return ingress
	}

	test := func(ingress *v1beta1.Ingress, ruleIndex int, pathOffset int, paths ...IngressPath) {
		for i, path := range paths {
			// The first two rules are dummies and should be skipped
			rulePath := ingress.Spec.Rules[ruleIndex].HTTP.Paths[i+pathOffset]
			suite.Equal(rulePath.Path, path.path)
			suite.Equal(rulePath.Backend.ServiceName, path.service)
			suite.Equal(rulePath.Backend.ServicePort.IntVal, int32(path.servicePort))
		}
	}

	// Add a single rule
	paths := []IngressPath{
		{
			path:        "/path/1/operation/123",
			service:     "service-1",
			servicePort: 9002,
		},
	}
	ingress := upsert(suite.ingress, suite.generateRulePath(paths[0]))
	test(ingress, suite.ingressHostRuleIndex, 1, paths...)

	// Add multiple rules
	multiplePaths := []IngressPath{
		{
			path:        "/path/2/operation/345",
			service:     "service-2",
			servicePort: 9003,
		},
		{
			path:        "/path/3/operation/456",
			service:     "service-3",
			servicePort: 9004,
		},
		{
			path:        "/path/4/operation/567",
			service:     "service-4",
			servicePort: 9005,
		},
	}
	ingress = upsert(
		suite.ingress,
		suite.generateRulePath(multiplePaths[0]),
		suite.generateRulePath(multiplePaths[1]),
		suite.generateRulePath(multiplePaths[2]),
	)
	test(ingress, suite.ingressHostRuleIndex, 2, multiplePaths...)

	// At this point, we have the following paths
	// 0. /dummy
	// 1. /path/1/operation/123
	// 2. /path/2/operation/345
	// 3. /path/3/operation/456
	// 4. /path/4/operation/567

	// Remove a non-existent path
	ingress = remove(ingress, "/non-existent")
	test(ingress, suite.ingressHostRuleIndex, 1, append(paths, multiplePaths...)...)

	// Remove a single non-dummy path
	ingress = remove(ingress, multiplePaths[0].path)
	test(ingress, suite.ingressHostRuleIndex, 1, paths[0], multiplePaths[2], multiplePaths[1])

	// Remove multiple non-dummy path
	ingress = remove(ingress, paths[0].path, "/non-existent", multiplePaths[1].path)
	test(ingress, suite.ingressHostRuleIndex, 1, multiplePaths[2])
}

func TestIsSubmissionDeadlineReached(t *testing.T) {
	t.Run("Should return false when submission deadline is not set", func(t *testing.T) {
		c := SubTCircuitRules{SubmissionDeadline: nil}
		assert.False(t, isSubmissionDeadlineReached(c))
	})

	t.Run("Should return false when submission deadline has not been reached", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		c := SubTCircuitRules{SubmissionDeadline: &deadline}
		assert.False(t, isSubmissionDeadlineReached(c))
	})

	t.Run("Should return true when submission deadline has been reached", func(t *testing.T) {
		deadline := time.Now()
		c := SubTCircuitRules{SubmissionDeadline: &deadline}
		assert.True(t, isSubmissionDeadlineReached(c))
	})
}
func TestIsCompetitionCircuit(t *testing.T) {
	assert.True(t, IsCompetitionCircuit("Tunnel Circuit"))
	assert.True(t, IsCompetitionCircuit("Urban Circuit"))
	assert.True(t, IsCompetitionCircuit("Cave Circuit"))

	assert.False(t, IsCompetitionCircuit("Tunnel Practice 1"))
	assert.False(t, IsCompetitionCircuit("Urban Practice 1"))
	assert.False(t, IsCompetitionCircuit("Cave Practice 1"))
}
