package simulations

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"net/http"
	"net/url"
)

// GetS3SimulationLogKey returns the key where logs for a given simulation are stored within a bucket.
func GetS3SimulationLogKey(dep *SimulationDeployment) string {
	groupID := *dep.GroupID
	ownerNameEscaped := url.PathEscape(*dep.Owner)
	key := fmt.Sprintf("/gz-logs/%s/%s/", ownerNameEscaped, groupID)

	return key
}

// UploadToS3Bucket uploads a file to a bucket in a certain key.
func UploadToS3Bucket(s3Svc s3iface.S3API, bucket *string, key *string, file []byte) (*s3.PutObjectOutput, error) {
	return s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:               bucket,
		Key:                  key,
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(file),
		ContentLength:        aws.Int64(int64(len(file))),
		ContentType:          aws.String(http.DetectContentType(file)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
}
