package cloud

// NewMock creates a new amazonWS instance with mocked components.
// EC2 -> AmazonEC2Mock
// S3 -> AmazonS3Mock
func NewMock() *amazonWS {
	ws := amazonWS{}
	ws.session = nil
	ws.ec2 = NewAmazonEC2Mock()
	ws.s3 = NewAmazonS3Mock()
	return &ws
}
