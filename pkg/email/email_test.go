package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailReturnsErrWhenRecipientsIsEmpty(t *testing.T) {
	s := NewEmailSender()

	err := s.Send(nil, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyRecipientList, err)
}

func TestEmailReturnsErrWhenSenderIsEmpty(t *testing.T) {
	s := NewEmailSender()

	err := s.Send([]string{"recipient@test.org"}, "", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrEmptySender, err)
}

func TestEmailReturnsErrWhenRecipientIsInvalid(t *testing.T) {
	s := NewEmailSender()

	err := s.Send([]string{"ThisIsNotAValidEmail"}, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidRecipient, err)
}

func TestEmailReturnsErrWhenSenderIsInvalid(t *testing.T) {
	s := NewEmailSender()

	err := s.Send([]string{"recipient@test.org"}, "InvalidSenderEmail", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidSender, err)
}
