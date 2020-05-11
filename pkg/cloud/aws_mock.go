package cloud

// NewMock creates a new AmazonWS instance with mocked components.
// EC2 -> AmazonEC2Mock
// S3 -> AmazonS3Mock
func NewMock() *AmazonWS {
	ws := AmazonWS{}
	ws.session = nil
	ws.EC2 = NewAmazonEC2Mock()
	ws.S3 = NewAmazonS3Mock()
	return &ws
}
