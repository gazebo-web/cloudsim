package email

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew_Null(t *testing.T) {
	email := New()

	assert.Zero(t, email.Sender())
	assert.Len(t, email.Recipients(), 0)
}

func TestNew_Empty(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "")

	email := New()

	assert.Zero(t, email.Sender())
	assert.Len(t, email.Recipients(), 0)
}

func TestNew_SetRecipient(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "sender@ignitionrobotics.org")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "recipient@ignitionrobotics.org")

	email := New()

	assert.Equal(t, "sender@ignitionrobotics.org", email.Sender())
	assert.Len(t, email.Recipients(), 1)
	assert.Equal(t, []string{"recipient@ignitionrobotics.org"}, email.Recipients())
}

func TestNew_SetRecipients(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "sender@ignitionrobotics.org")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "recipient@ignitionrobotics.org,another@ignitionrobotics.org,example@ignitionrobotics.org")

	email := New()

	assert.Equal(t, email.Sender(), "sender@ignitionrobotics.org")
	assert.Len(t, email.Recipients(), 3)
	assert.Equal(t, []string{"recipient@ignitionrobotics.org", "another@ignitionrobotics.org", "example@ignitionrobotics.org"}, email.Recipients())
}
