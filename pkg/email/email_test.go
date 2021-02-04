package email

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailReturnsErrWhenRecipientsIsEmpty(t *testing.T) {
	s := NewEmailSender(nil)

	err := s.Send(nil, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyRecipientList, err)
}

func TestEmailReturnsErrWhenSenderIsEmpty(t *testing.T) {
	s := NewEmailSender(nil)

	err := s.Send([]string{"recipient@test.org"}, "", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrEmptySender, err)
}

func TestEmailReturnsErrWhenRecipientIsInvalid(t *testing.T) {
	s := NewEmailSender(nil)

	err := s.Send([]string{"ThisIsNotAValidEmail"}, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidRecipient, err)
}

func TestEmailReturnsErrWhenSenderIsInvalid(t *testing.T) {
	s := NewEmailSender(nil)

	err := s.Send([]string{"recipient@test.org"}, "InvalidSenderEmail", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidSender, err)
}

// fakeSender fakes the sesiface.SESAPI interface.
type fakeSender struct {
	returnError bool
	sesiface.SESAPI
}

// SendEmail mocks the SendEmail method from the sesiface.SESAPI.
func (s *fakeSender) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	if s.returnError {
		return nil, errors.New("fake error")
	}
	return nil, nil
}
