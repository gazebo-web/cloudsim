package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/calculator"
)

const (
	// KindMachines is used to identify the kind of calculator.Resource for EC2 machines.
	KindMachines = "AmazonEC2"
	// KindStorage is used to identify the kind of calculator.Resource for S3 storage.
	KindStorage = "AmazonS3"
)

// PriceParser parses a given AWS price result from the Pricing API and returns a calculator.Rate value.
// Different implementations depending on the response of the Pricing API can be added using this function signature.
// EC2 example: ParseEC2
type PriceParser func(price aws.JSONValue) (calculator.Rate, error)

// costCalculator provides an AWS calculator for different services.
type costCalculator struct {
	// API holds a reference to a pricingiface.PricingAPI implementation.
	API pricingiface.PricingAPI
	// priceParser holds a reference to a specific price parser implementation such as ParseEC2.
	priceParser PriceParser
	// resource describes the name of the service that is being used to calculate costs from.
	resource string
}

// CalculateCost calculates the cost of an arbitrary set of AWS resources.
func (c *costCalculator) CalculateCost(resources []calculator.Resource) (calculator.Rate, error) {
	rates := make([]calculator.Rate, len(resources))
	for i, res := range resources {
		rate, err := c.calculateRate(res)
		if err != nil {
			return calculator.Rate{}, err
		}
		rates[i] = rate
	}
	return calculator.AggregateRates(rates), nil
}

// calculateRate calculates the rate of a given resource.
func (c *costCalculator) calculateRate(res calculator.Resource) (calculator.Rate, error) {
	filters := c.convertResourceToFilters(res)
	filters = c.appendServiceCodeFilter(filters)
	list, err := c.API.GetProducts(&pricing.GetProductsInput{
		FormatVersion: aws.String("aws_v1"),
		MaxResults:    aws.Int64(1),
		ServiceCode:   aws.String(c.resource),
		Filters:       filters,
	})
	if err != nil {
		return calculator.Rate{}, err
	}
	if len(list.PriceList) == 0 {
		return calculator.Rate{}, errors.New("product not found")
	}
	rate, err := c.priceParser(list.PriceList[0])
	if err != nil {
		return calculator.Rate{}, err
	}
	return rate, nil
}

// convertResourceToFilters converts the given resource into a set of filters for the Pricing API.
func (c *costCalculator) convertResourceToFilters(resource calculator.Resource) []*pricing.Filter {
	filters := make([]*pricing.Filter, 0, len(resource.Values))
	for k, v := range resource.Values {
		value, ok := v.(string)
		if !ok {
			continue
		}
		filters = append(filters, &pricing.Filter{
			Field: aws.String(k),
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Value: &value,
		})
	}
	return filters
}

// appendServiceCodeFilter adds a service code filter to the given list of filters.
func (c *costCalculator) appendServiceCodeFilter(filters []*pricing.Filter) []*pricing.Filter {
	return append(filters, &pricing.Filter{
		Field: aws.String("ServiceCode"),
		Type:  aws.String(pricing.FilterTypeTermMatch),
		Value: &c.resource,
	})
}

// NewCostCalculator initializes a new cost calculator for a certain AWS resource.
func NewCostCalculator(api pricingiface.PricingAPI, priceParser PriceParser, resource string) calculator.CostCalculator {
	return &costCalculator{
		API:         api,
		priceParser: priceParser,
		resource:    resource,
	}
}
