package email

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Send sends an email to a specific recipient. If the recipient is nil,
// then the default recipient defined in the IGN_FLAGS_EMAIL_TO env var will be
// used.
func (e *Email) Send(recipient *[]string, sender *string, subject string, templateFilename string, templateData interface{}) *ign.ErrMsg {
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