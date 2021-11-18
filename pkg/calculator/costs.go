package calculator

import "time"

// Rate is the rate at what a certain resource is charged.
type Rate struct {
	// Amount is the money a resource costs in the minimum Currency value (e.g. cents for USD) at a certain Frequency.
	Amount uint

	// Currency is the ISO 4217 currency code in lowercase format.
	Currency string

	// Frequency is the frequency at which a resource gets charged.
	Frequency time.Duration
}

// Aggregate merges the current rate with the rate given as an argument, and returns the sum of both values
// expressed in hours.
func (r Rate) Aggregate(rate Rate) Rate {
	r = transformRate(r, time.Hour)
	rate = transformRate(rate, time.Hour)
	return Rate{
		Amount:    r.Amount + rate.Amount,
		Currency:  rate.Currency,
		Frequency: time.Hour,
	}
}

// transformRate transforms the given rate to the given frequency.
// If the frequency of the given rate is smaller than the one passed as the freq argument
// it will return the representation of a rate using the given frequency.
// Otherwise, it will return the current rate.
func transformRate(rate Rate, freq time.Duration) Rate {
	f := int64(1)
	if freq.Milliseconds() > rate.Frequency.Milliseconds() && rate.Frequency.Milliseconds() > 0 {
		f = freq.Milliseconds() / rate.Frequency.Milliseconds()
		rate.Frequency = freq
	}
	rate.Amount = rate.Amount * uint(f)
	return rate
}

// SumRates sums up the given rates and returns the representation in hours.
func SumRates(rates []Rate) Rate {
	var out Rate
	for _, r := range rates {
		out = out.Aggregate(r)
	}
	return out
}

// RateAggregator holds a method to merge two rates together.
type RateAggregator interface {
	// Aggregate merges the current rate with the rate given as an argument, and returns the sum of both values
	// expressed in hours.
	Aggregate(rate Rate) Rate
}

// Resource groups a set of fields from a resource consumed by cloudsim. It's used to calculate the cost at which
// the resource should be charged for.
type Resource struct {
	// Kind represents what kind of resource is being used to calculate its Rate.
	Kind string

	// Values contains a set of arbitrary values used to calculate the Rate of this resource.
	Values map[string]interface{}
}

// CostCalculator holds a method to calculate the cost of a group of resources.
type CostCalculator interface {
	// CalculateCost calculates the Rate at which a group of resources should be charged for.
	CalculateCost(resources []Resource) (Rate, error)
}
