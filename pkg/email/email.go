package email

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Config represents a set of options to configure the Email service.
type Email struct {
	DefaultEmailRecipients []string
	DefaultEmailSender string
}

// TODO: Find a better name. This should be its own service and not only a configuration parser.
// New returns a new Email configuration.
func New() *Email {
	email := Email{}
	email.DefaultEmailRecipients = tools.EnvVarToSlice("IGN_DEFAULT_EMAIL_RECIPIENT")
	email.DefaultEmailSender, _ = ign.ReadEnvVar("IGN_DEFAULT_EMAIL_SENDER")
	return &email
}