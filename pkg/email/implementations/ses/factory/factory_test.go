package factory

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestS3FactorySuite(t *testing.T) {
	suite.Run(t, new(testSESFactorySuite))
}

type testSESFactorySuite struct {
	suite.Suite
}

func (s *testSESFactorySuite) TestInitializeAPIDependencyIsNil() {
	config := Config{
		Region: "us-east-1",
	}
	dependencies := Dependencies{}
	s.Nil(initializeAPI(&config, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testSESFactorySuite) TestInitializeAPIDependencyIsNotNil() {
	// Prepare dependencies
	sesAPI := struct {
		sesiface.SESAPI
	}{}
	dependencies := Dependencies{
		API: sesAPI,
	}

	s.Nil(initializeAPI(nil, &dependencies))
	s.Exactly(sesAPI, dependencies.API)
}
