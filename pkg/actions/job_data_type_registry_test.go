package actions

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistryRegisterAndGetType(t *testing.T) {
	r := jobDataTypeRegistry

	type TestInput struct{}
	input := TestInput{}

	register := func(value interface{}) {
		r.register(GetJobDataType(value))
		out, err := r.getType(GetJobDataTypeName(value))
		require.NoError(t, err)
		if value != nil {
			require.NotPanics(t, func() {
				_ = out
			})
		} else {
			require.Nil(t, out)
		}
	}

	// Concrete type
	register(input)

	// Concrete type pointer
	register(&input)

	// Nil
	register(nil)
}

func TestRegistryGetTypeNotFound(t *testing.T) {
	r := jobDataTypeRegistry

	_, err := r.getType("fail")
	require.Error(t, err)
}
