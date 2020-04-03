package cloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// AmazonWS represents an Amazon Web Service client instance.
type AmazonWS struct {
	// From aws go documentation:
	// Sessions should be cached when possible, because creating a new Session
	// will load all configuration values from the environment, and config files
	// each time the Session is created. Sharing the Session value across all of
	// your service clients will ensure the configuration is loaded the fewest
	// number of times possible.
	session *session.Session
	EC2 *AmazonEC2
	S3 *AmazonS3
}

// New returns a new Amazon Web Service client.
func New() *AmazonWS {
	ws := AmazonWS{}
	ws.session = session.Must(session.NewSession())
	ws.EC2 = NewAmazonEC2(ws.session)
	ws.S3 = NewAmazonS3(ws.session)
	return &ws
}
