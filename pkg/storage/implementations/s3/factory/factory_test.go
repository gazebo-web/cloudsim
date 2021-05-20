package factory

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestS3FactorySuite(t *testing.T) {
	suite.Run(t, new(testS3FactorySuite))
}

type testS3FactorySuite struct {
	suite.Suite
}

func (s *testS3FactorySuite) TestInitializeAPIDependencyIsNil() {
	config := Config{
		Region: "us-east-1",
	}
	dependencies := Dependencies{}
	s.Nil(initializeAPI(&config, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testS3FactorySuite) TestInitializeAPIDependencyIsNotNil() {
	// Prepare dependencies
	s3API := struct {
		s3iface.S3API
	}{}
	dependencies := Dependencies{
		API: s3API,
	}

	s.Nil(initializeAPI(nil, &dependencies))
	s.Exactly(s3API, dependencies.API)
}
