package simulations

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceTag(t *testing.T) {
	input := &ec2.RunInstancesInput{
		TagSpecifications: []*ec2.TagSpecification{
			{
				Tags: []*ec2.Tag{
					{
						Key:   sptr("key_1"),
						Value: sptr("value_1"),
					},
					{
						Key:   sptr("key_2"),
						Value: sptr("value_2"),
					},
					{
						Key:   sptr("key_3"),
						Value: sptr("value_3"),
					},
				},
			},
		},
	}

	test := func(values ...string) {
		tags := input.TagSpecifications[0].Tags
		assert.Equal(t, len(tags), len(values))
		for i := range tags {
			assert.Equal(t, *tags[i].Value, values[i])
		}
	}

	// Invalid key
	tag := &ec2.Tag{
		Key:   sptr("invalid"),
		Value: sptr("test"),
	}
	replaceTag(input, tag)
	test("value_1", "value_2", "value_3")

	// Replace values 2 and 3
	tag2 := &ec2.Tag{
		Key:   sptr("key_2"),
		Value: sptr("test_2"),
	}
	tag3 := &ec2.Tag{
		Key:   sptr("key_3"),
		Value: sptr("test_3"),
	}
	replaceTag(input, tag2, tag3)
	test("value_1", "test_2", "test_3")
}
