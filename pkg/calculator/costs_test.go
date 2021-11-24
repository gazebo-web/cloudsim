package calculator

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSumRates(t *testing.T) {
	r1 := Rate{
		Amount:    100, // 100 in seconds -> 360000 in hours
		Currency:  "usd",
		Frequency: time.Second,
	}

	r2 := Rate{
		Amount:    100, // 100 in minutes -> 6000 in hours
		Currency:  "usd",
		Frequency: time.Minute,
	}

	r3 := Rate{
		Amount:    100, // 100 in hours -> 100 in hours.
		Currency:  "usd",
		Frequency: time.Hour,
	}

	out := AggregateRates([]Rate{r1, r2, r3})

	assert.Equal(t, "usd", out.Currency)
	assert.Equal(t, time.Hour, out.Frequency)
	assert.Equal(t, uint(360000+6000+100), out.Amount)
}

func TestTransformRate(t *testing.T) {
	// If the frequency of the current rate is smaller than the one passed in transformRate,
	// it should return the representation of a rate using the given frequency.

	// Seconds to Minutes should return Rate in minutes.
	in := Rate{
		Amount:    100, // 1 usd
		Currency:  "usd",
		Frequency: time.Second,
	}
	out := transformRate(in, time.Minute)

	assert.Equal(t, uint(6000), out.Amount)
	assert.Equal(t, time.Minute, out.Frequency)

	// If the frequency of the current rate is bigger than the one passed in transformRate,
	// it should default to the current rate's frequency, as an invalid conversion could happen when passing a lower rate.

	// Hour to Seconds should return Rate in hours.
	in = Rate{
		Amount:    2000, // 20 usd
		Currency:  "usd",
		Frequency: time.Hour,
	}
	out = transformRate(in, time.Second)
	assert.Equal(t, uint(2000), out.Amount)
	assert.Equal(t, time.Hour, out.Frequency)
}
