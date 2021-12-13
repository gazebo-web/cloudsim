package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestPriceParserEC2(t *testing.T) {
	b, err := ioutil.ReadFile("./price_ec2.json")
	require.NoError(t, err)

	var value aws.JSONValue
	require.NoError(t, json.Unmarshal(b, &value))

	rate, err := ParseEC2(value)
	require.NoError(t, err)
	assert.Equal(t, time.Hour, rate.Frequency)
	assert.Equal(t, "usd", rate.Currency)
	assert.Equal(t, uint(114), rate.Amount)
}

func TestNormalizeAmount(t *testing.T) {
	assert.Equal(t, uint(114), parseAmount(1.14))
	assert.Equal(t, uint(235), parseAmount(2.35))
	assert.Equal(t, uint(1514), parseAmount(15.14))
	assert.Equal(t, uint(70), parseAmount(0.7))
	assert.Equal(t, uint(104), parseAmount(1.04))
	assert.Equal(t, uint(30), parseAmount(0.3))
}
