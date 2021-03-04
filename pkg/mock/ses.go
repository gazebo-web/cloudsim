package mock

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// ec2api is an ec2iface.EC2API implementation.
type sesAPI struct {
	sesiface.SESAPI
}

// SendEmail mocks the SES API SendEmail method.
func (s *sesAPI) SendEmail(*ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return nil, nil
}

// NewEC2 initializes a new ec2iface.EC2API implementation.
func NewSES() sesiface.SESAPI {
	return &sesAPI{}
}
