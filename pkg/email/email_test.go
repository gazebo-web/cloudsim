package email

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew_Null(t *testing.T) {
	email := New()

	assert.Zero(t, email.DefaultEmailSender)
	assert.Len(t, email.DefaultEmailRecipients, 0)
}

func TestNew_Empty(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "")

	email := New()

	assert.Zero(t, email.DefaultEmailSender)
	assert.Len(t, email.DefaultEmailRecipients, 0)
}

func TestNew_SetRecipient(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "sender@ignitionrobotics.org")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "recipient@ignitionrobotics.org")

	email := New()

	assert.Equal(t, email.DefaultEmailSender, "sender@ignitionrobotics.org")
	assert.Len(t, email.DefaultEmailRecipients, 1)
	assert.Equal(t, []string{"recipient@ignitionrobotics.org"}, email.DefaultEmailRecipients)
}

func TestNew_SetRecipients(t *testing.T) {
	os.Setenv("IGN_DEFAULT_EMAIL_SENDER", "sender@ignitionrobotics.org")
	os.Setenv("IGN_DEFAULT_EMAIL_RECIPIENT", "recipient@ignitionrobotics.org,another@ignitionrobotics.org,example@ignitionrobotics.org")

	email := New()

	assert.Equal(t, email.DefaultEmailSender, "sender@ignitionrobotics.org")
	assert.Len(t, email.DefaultEmailRecipients, 3)
	assert.Equal(t, []string{"recipient@ignitionrobotics.org", "another@ignitionrobotics.org", "example@ignitionrobotics.org"}, email.DefaultEmailRecipients)
}