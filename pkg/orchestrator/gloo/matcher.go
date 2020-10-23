package gloo

import "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

// Matcher groups a set of methods to represent a route matcher.
type Matcher interface {
	// Value returns the value of the matcher.
	Value() string

	// ToExact returns an exact matcher.
	ToExact() matchers.Matcher

	// ToPrefix returns a prefix matcher.
	ToPrefix() matchers.Matcher

	// ToRegex returns a regex matcher.
	ToRegex() matchers.Matcher
}

// match is a matcher representation.
type match struct {
	// expressions defines the matcher value.
	expression string
}

// Value returns the value of the expression associated with this matcher.
func (m match) Value() string {
	return m.expression
}

// ToExact converts the underlying expression into a matchers.Matcher of type matchers.Matcher_Exact.
func (m match) ToExact() matchers.Matcher {
	return matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Exact{
			Exact: m.Value(),
		},
	}
}

// ToPrefix converts the underlying expression into a matchers.Matcher of type matchers.Matcher_Prefix.
func (m match) ToPrefix() matchers.Matcher {
	return matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Prefix{
			Prefix: m.Value(),
		},
	}
}

// ToRegex converts the underlying expression into a matchers.Matcher of type matchers.Matcher_Regex.
func (m match) ToRegex() matchers.Matcher {
	return matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: m.Value(),
		},
	}
}

// NewMatcher initializes a new matcher with the given expression.
func NewMatcher(expr string) Matcher {
	return &match{
		expression: expr,
	}
}
