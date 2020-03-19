package email

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Email struct {
	DefaultEmailRecipients []string
	DefaultEmailSender string
}

func New() Email {
	email := Email{}
	email.DefaultEmailRecipients = tools.EnvVarToSlice("IGN_DEFAULT_EMAIL_RECIPIENT")
	email.DefaultEmailSender, _ = ign.ReadEnvVar("IGN_DEFAULT_EMAIL_SENDER")
	return email
}