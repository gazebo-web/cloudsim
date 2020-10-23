package gloo

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	var m Matcher

	expr := "test"

	m = NewMatcher(expr)

	assert.Equal(t, expr, m.Value())

	exact := matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Exact{
			Exact: expr,
		},
	}

	assert.True(t, exact.Equal(m.ToExact()))

	prefix := matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Prefix{
			Prefix: expr,
		},
	}

	assert.True(t, prefix.Equal(m.ToPrefix()))

	regex := matchers.Matcher{
		PathSpecifier: &matchers.Matcher_Regex{
			Regex: expr,
		},
	}

	assert.True(t, regex.Equal(m.ToRegex()))
}
