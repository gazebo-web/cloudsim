package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestPriceParserEC2(t *testing.T) {
	b, err := os.ReadFile("./price_ec2.json")
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
	assert.Equal(t, uint(114), normalizeAmount(1.14))
	assert.Equal(t, uint(235), normalizeAmount(2.35))
	assert.Equal(t, uint(1514), normalizeAmount(15.14))
	assert.Equal(t, uint(70), normalizeAmount(0.7))
	assert.Equal(t, uint(104), normalizeAmount(1.04))
	assert.Equal(t, uint(30), normalizeAmount(0.3))
}
