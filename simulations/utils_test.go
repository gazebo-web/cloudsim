package simulations

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Since tokens are in hex, they should be twice as long as the input byte size.
	test := func(size *int) {
		if size == nil {
			size = intptr(32)
		}
		tokenSize := *size * 2
		token, err := generateToken(size)
		assert.NoError(t, err)
		assert.Len(t, token, tokenSize)
		assert.Regexp(t, fmt.Sprintf("[a-f0-9]{%d}", tokenSize), token)
	}

	// nil size
	test(nil)

	// Size 0
	test(intptr(0))

	// Explicit size
	test(intptr(8))
}
