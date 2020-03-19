package auth0

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


func TestNew_Null(t *testing.T) {
	auth := New()
	assert.Zero(t, auth.PublicKey)
}

func TestNew_Empty(t *testing.T) {
	os.Setenv("AUTH0_RSA256_PUBLIC_KEY", "")

	auth := New()
	assert.Zero(t, auth.PublicKey)
}

func TestNew_Set(t *testing.T) {
	os.Setenv("AUTH0_RSA256_PUBLIC_KEY", "auth-0-public-key-value")

	auth := New()

	assert.Equal(t, "auth-0-public-key-value", auth.PublicKey)
}
