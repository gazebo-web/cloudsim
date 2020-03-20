package email

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Config struct {
	DefaultEmailRecipients []string
	DefaultEmailSender string
}

func New() Config {
	email := Config{}
	email.DefaultEmailRecipients = tools.EnvVarToSlice("IGN_DEFAULT_EMAIL_RECIPIENT")
	email.DefaultEmailSender, _ = ign.ReadEnvVar("IGN_DEFAULT_EMAIL_SENDER")
	return email
}