package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/calculator"
	"os"
	"testing"
)

type pricingMock struct {
	pricingiface.PricingAPI
	T                *testing.T
	ExpectedResource calculator.Resource
}

func (api *pricingMock) GetProducts(input *pricing.GetProductsInput) (*pricing.GetProductsOutput, error) {
	require.NotNil(api.T, input)
	require.NotNil(api.T, input.FormatVersion)
	require.NotNil(api.T, input.ServiceCode)
	require.NotNil(api.T, input.MaxResults)

	assert.Equal(api.T, "aws_v1", *input.FormatVersion)
	assert.Equal(api.T, "AmazonEC2", *input.ServiceCode)
	assert.Equal(api.T, int64(1), *input.MaxResults)

	require.NotEmpty(api.T, input.Filters)
	for _, r := range input.Filters {
		require.NotNil(api.T, r.Field)
		require.NotNil(api.T, r.Type)
		require.NotNil(api.T, r.Value)

		assert.Equal(api.T, pricing.FilterTypeTermMatch, *r.Type)
		v, ok := api.ExpectedResource.Values[*r.Field]
		assert.True(api.T, ok)
		assert.Equal(api.T, v, *r.Value)
	}

	b, err := os.ReadFile("./price_ec2.json")
	require.NoError(api.T, err)

	var value aws.JSONValue
	require.NoError(api.T, json.Unmarshal(b, &value))

	return &pricing.GetProductsOutput{
		FormatVersion: nil,
		NextToken:     nil,
		PriceList:     []aws.JSONValue{value},
	}, nil
}

func TestCalculator(t *testing.T) {
	res := calculator.Resource{
		Values: map[string]interface{}{
			"ServiceCode":     "AmazonEC2",
			"instanceType":    "g3.4xlarge",
			"marketoption":    "OnDemand",
			"operatingSystem": "Linux",
			"regionCode":      "us-east-1",
			"tenancy":         "Shared",
			"capacitystatus":  "Used",
		},
	}

	c := newCostCalculator(&pricingMock{
		T:                t,
		ExpectedResource: res,
	}, ParseEC2, KindMachines)

	rate, err := c.CalculateCost([]calculator.Resource{res})
	require.NoError(t, err)

	assert.Equal(t, uint(114), rate.Amount)
}
