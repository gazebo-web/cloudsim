package ses

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
)

// NewAPI returns an AWS Simple Email Service (SES) client from the given config provider.
func NewAPI(config client.ConfigProvider) sesiface.SESAPI {
	return ses.New(config)
}

// email implements the Sender interface using AWS SES api.
type sesEmail struct {
	API    sesiface.SESAPI
	Logger ign.Logger
}

// Send sends an email to the given recipients from the given sender.
// A template will be parsed with the given data in order to fill the email's body.
// It returns an error when validation fails or sending an email fails.
func (e *sesEmail) Send(recipients []string, sender, subject, template string, data interface{}) error {
	if len(recipients) == 0 {
		return email.ErrEmptyRecipientList
	}
	if len(sender) == 0 {
		return email.ErrEmptySender
	}
	for _, r := range recipients {
		if ok := email.ValidateEmail(r); !ok {
			return email.ErrInvalidRecipient
		}
	}
	if ok := email.ValidateEmail(sender); !ok {
		return email.ErrInvalidSender
	}
	if data == nil {
		return email.ErrInvalidData
	}

	content, err := ign.ParseHTMLTemplate(template, data)
	if err != nil {
		return err
	}

	err = e.send(sender, recipients, subject, content)

	return nil
}

// send attempts to send an email using the AWS SES service.
func (e *sesEmail) send(sender string, recipients []string, subject string, content string) error {
	input := ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(recipients),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(email.CharSet),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(email.CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the sesEmail.
	_, err := e.API.SendEmail(&input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			var code string
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				code = ses.ErrCodeMessageRejected
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				code = ses.ErrCodeMailFromDomainNotVerifiedException
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				code = ses.ErrCodeConfigurationSetDoesNotExistException
			default:
				code = "Unknown AWS SES error"
			}
			return errors.New(code + " " + aerr.Error())
		}
		return errors.New(err.Error())
	}

	return nil
}

// NewEmailSender returns a email.Sender implementation.
func NewEmailSender(api sesiface.SESAPI, logger ign.Logger) email.Sender {
	return &sesEmail{
		API:    api,
		Logger: logger,
	}
}
