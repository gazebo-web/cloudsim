package ses

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"testing"
)

func TestEmailReturnsErrWhenRecipientsIsEmpty(t *testing.T) {
	s := NewEmailSender(nil, nil)

	err := s.Send(nil, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, email.ErrEmptyRecipientList, err)
}

func TestEmailReturnsErrWhenSenderIsEmpty(t *testing.T) {
	s := NewEmailSender(nil, nil)

	err := s.Send([]string{"recipient@test.org"}, "", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, email.ErrEmptySender, err)
}

func TestEmailReturnsErrWhenRecipientIsInvalid(t *testing.T) {
	s := NewEmailSender(nil, nil)

	err := s.Send([]string{"ThisIsNotAValidEmail"}, "example@test.org", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, email.ErrInvalidRecipient, err)
}

func TestEmailReturnsErrWhenSenderIsInvalid(t *testing.T) {
	s := NewEmailSender(nil, nil)

	err := s.Send([]string{"recipient@test.org"}, "InvalidSenderEmail", "Some test", "test.template", nil)
	assert.Error(t, err)
	assert.Equal(t, email.ErrInvalidSender, err)
}

func TestEmailReturnsErrWhenInvalidPath(t *testing.T) {
	s := NewEmailSender(nil, nil)

	err := s.Send([]string{"recipient@test.org"}, "example@test.org", "Some test", "test", nil)
	assert.Error(t, err)
}

func TestEmailReturnsErrWhenDataIsNil(t *testing.T) {
	s := NewEmailSender(&fakeSender{}, nil)
	err := s.Send([]string{"recipient@test.org"}, "example@test.org", "Some test", "template.gohtml", nil)
	assert.Error(t, err)
	assert.Equal(t, email.ErrInvalidData, err)
}

func TestEmailSendingSuccess(t *testing.T) {
	fake := fakeSender{}
	s := NewEmailSender(&fake, nil)
	err := s.Send([]string{"recipient@test.org"}, "example@test.org", "Some test", "template.gohtml", struct {
		Test string
	}{
		Test: "Hello there!",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, fake.Called)
}

// fakeSender fakes the sesiface.SESAPI interface.
type fakeSender struct {
	returnError bool
	Called      int
	sesiface.SESAPI
}

// SendEmail mocks the SendEmail method from the sesiface.SESAPI.
func (s *fakeSender) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	s.Called++
	if s.returnError {
		return nil, errors.New("fake error")
	}
	return nil, nil
}
