package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/itchyny/gojq"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/calculator"
	"math"
	"strconv"
	"strings"
	"time"
)

// priceEC2 is a data structure used as a helper to represent
type priceEC2 struct {
	// Frequency is the frequency at which a certain EC2 machines gets charged at.
	// Values: Hrs
	Frequency string `json:"frequency"`

	// Amounts groups the currencies and the respective amounts of money that a certain EC2 machines can be charged for.
	Amounts []struct {
		// Currency holds the currency code
		Currency string `json:"currency"`
		// Amount is the amount of money to charge.
		Amount string `json:"amount"`
	} `json:"amounts"`
}

// ParseEC2 is a PriceParser func used for parsing EC2 pricing. It reads the given product definition and returns
// a rate at which the given product should be charged.
// If multiple currencies are present, it returns the first one it finds.
func ParseEC2(product aws.JSONValue) (calculator.Rate, error) {
	q, err := gojq.Parse("{frequency: .terms.OnDemand[].priceDimensions[].unit, amounts: .terms.OnDemand[].priceDimensions[].pricePerUnit | to_entries | map_values({ currency: .key, amount: .value }) }")
	if err != nil {
		return calculator.Rate{}, err
	}

	iter := q.Run(map[string]interface{}(product))
	v, ok := iter.Next()
	if !ok {
		return calculator.Rate{}, errors.New("failed to parse JSON using jq")
	}

	p, err := convertMapToPriceEC2(v)
	if err != nil {
		return calculator.Rate{}, err
	}

	amount, err := strconv.ParseFloat(p.Amounts[0].Amount, 64)
	if err != nil {
		return calculator.Rate{}, err
	}

	return calculator.Rate{
		Amount:    normalizeAmount(amount),
		Currency:  normalizeCurrency(p.Amounts[0].Currency),
		Frequency: normalizeFrequency(p.Frequency),
	}, nil
}

// normalizeFrequency converts the given unit of time from AWS to time.Duration.
// It defaults to time.Hour.
func normalizeFrequency(timeUnit string) time.Duration {
	switch timeUnit {
	case "Hrs":
		return time.Hour
	}
	return time.Hour
}

// normalizeCurrency converts the given currency into a ISO 4217 currency code in lowercase format.
func normalizeCurrency(currency string) string {
	return strings.ToLower(currency)
}

// normalizeAmount converts the given amount to cents in USD.
func normalizeAmount(amount float64) uint {
	return uint(math.Round(amount * 100))
}

// convertMapToPriceEC2 decodes the given map as a priceEC2 structure.
func convertMapToPriceEC2(v interface{}) (priceEC2, error) {
	var p priceEC2
	if err := mapstructure.Decode(v, &p); err != nil {
		return priceEC2{}, err
	}
	return p, nil
}
