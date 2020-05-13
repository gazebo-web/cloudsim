package cloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type AmazonWS interface {
	Session() *session.Session
	EC2() AmazonEC2
	S3() AmazonS3
}

// amazonWS represents an Amazon Web Service client instance.
type amazonWS struct {
	// From aws go documentation:
	// Sessions should be cached when possible, because creating a new Session
	// will load all configuration values from the environment, and config files
	// each time the Session is created. Sharing the Session value across all of
	// your service clients will ensure the configuration is loaded the fewest
	// number of times possible.
	session *session.Session
	ec2     AmazonEC2
	s3      AmazonS3
}

func (ws *amazonWS) Session() *session.Session {
	return ws.session
}

func (ws *amazonWS) EC2() AmazonEC2 {
	return ws.ec2
}

func (ws *amazonWS) S3() AmazonS3 {
	return ws.s3
}

// New returns a new Amazon Web Service client wrapper.
func New() AmazonWS {
	ws := amazonWS{}
	ws.session = session.Must(session.NewSession())
	ws.ec2 = NewAmazonEC2(ws.session)
	ws.s3 = NewAmazonS3(ws.session)
	return &ws
}
