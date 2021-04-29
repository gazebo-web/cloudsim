package email

import (
	"errors"
	"regexp"
)

// CharSet is the default email encoding.
const CharSet = "UTF-8"

var (
	// ErrEmptyRecipientList is returned when an empty recipient list is passed to the Sender.Send method.
	ErrEmptyRecipientList = errors.New("empty recipient list")
	// ErrEmptySender is returned when an empty sender email address is passed to the Sender.Send method.
	ErrEmptySender = errors.New("empty sender")
	// ErrInvalidSender is returned when an invalid email address is passed to the Sender.Send method.
	ErrInvalidSender = errors.New("invalid sender")
	// ErrInvalidRecipient is returned when an invalid email is passed in the list of recipients to the Sender.Send method.
	ErrInvalidRecipient = errors.New("invalid recipient")
	// ErrInvalidData is returned when an invalid data is passed to the Sender.Send method.
	ErrInvalidData = errors.New("invalid data")
)

// Sender has a method to send emails.
type Sender interface {
	Send(recipients []string, sender, subject, template string, data interface{}) error
}

// ValidateEmail validates the given email. It returns false if the validation fails.
func ValidateEmail(email string) bool {
	exp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return exp.MatchString(email)
}
