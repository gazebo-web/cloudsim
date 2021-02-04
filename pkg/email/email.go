package email

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"regexp"
)

const charSet = "UTF-8"

var (
	// ErrEmptyRecipientList is returned when an empty recipient list is passed to the Sender.Send method.
	ErrEmptyRecipientList = errors.New("empty recipient list")
	// ErrEmptySender is returned when an empty sender email address is passed to the Sender.Send method.
	ErrEmptySender = errors.New("empty sender")
	// ErrInvalidSender is returned when an invalid email address is passed to the Sender.Send method.
	ErrInvalidSender = errors.New("invalid sender")
	// ErrInvalidRecipient is returned when an invalid email is passed in the list of recipients to the Sender.Send method.
	ErrInvalidRecipient = errors.New("invalid recipient")
)

// Sender has a method to send emails.
type Sender interface {
	Send(recipients []string, sender, subject, template string, data interface{}) error
}

type email struct {
	API sesiface.SESAPI
}

// Send sends an email to the given recipients from the given sender.
// A template will be parsed with the given data in order to fill the email's body.
// It returns an error when validation fails or sending an email fails.
func (e *email) Send(recipients []string, sender, subject, template string, data interface{}) error {
	if len(recipients) == 0 {
		return ErrEmptyRecipientList
	}
	if len(sender) == 0 {
		return ErrEmptySender
	}
	for _, r := range recipients {
		if ok := ValidateEmail(r); !ok {
			return ErrInvalidRecipient
		}
	}
	if ok := ValidateEmail(sender); !ok {
		return ErrInvalidSender
	}

	content, err := ign.ParseHTMLTemplate(template, data)
	if err != nil {
		return err
	}

	err = e.send(sender, recipients, subject, content)

	return nil
}

// send attempts to send an email using AWS SES service.
func (e *email) send(sender string, recipients []string, subject string, content string) error {
	input := ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(recipients),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(content),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
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

// NewEmailSender returns a Sender implementation.
func NewEmailSender(api sesiface.SESAPI) Sender {
	return &email{
		API: api,
	}
}

// ValidateEmail validates the given email. It returns false if the validation fails.
func ValidateEmail(email string) bool {
	exp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return exp.MatchString(email)
}
