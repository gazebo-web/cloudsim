package email

import (
	"errors"
	"regexp"
)

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

type email struct{}

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
	return nil
}

// NewEmailSender returns a Sender implementation.
func NewEmailSender() Sender {
	return &email{}
}

// ValidateEmail validates the given email. It returns false if the validation fails.
func ValidateEmail(email string) bool {
	exp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return exp.MatchString(email)
}
