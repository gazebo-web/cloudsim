package cloud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	aws := New()
	assert.NotNil(t, aws.session)
	assert.NotNil(t, aws.EC2)
	assert.NotNil(t, aws.S3)
}
