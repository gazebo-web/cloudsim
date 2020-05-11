package email

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Email interface {
	Send(recipient *[]string, sender *string, subject string, templateFilename string, templateData interface{}) *ign.ErrMsg
}

// Config represents a set of options to configure the email service.
type email struct {
	DefaultEmailRecipients []string `env:"IGN_DEFAULT_EMAIL_RECIPIENT" envSeparator:","`
	DefaultEmailSender     string `env:"IGN_DEFAULT_EMAIL_SENDER"`
}

// New returns a new email service.
func New() Email {
	email := email{}
	_ = env.Parse(&email)
	return &email
}

// Send sends an email to a specific recipient. If the recipient is nil,
// then the default recipient defined in the IGN_FLAGS_EMAIL_TO env var will be
// used.
func (e *email) Send(recipient *[]string, sender *string, subject string, templateFilename string, templateData interface{}) *ign.ErrMsg {
	if recipient == nil {
		recipient = &e.DefaultEmailRecipients
	}
	if sender == nil {
		sender = &e.DefaultEmailSender
	}
	// If the sender or recipient are not defined, then don't send the email
	if (recipient != nil && len(*recipient) == 0) || (sender != nil && *sender == "") {
		return nil
	}

	// Prepare the template
	content, err := ign.ParseHTMLTemplate(templateFilename, templateData)
	if err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// Send the email
	for _, r := range *recipient {
		err = ign.SendEmail(*sender, r, subject, content)
		if err != nil {
			return ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
		}
	}
	return nil
}
