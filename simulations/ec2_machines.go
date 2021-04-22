package simulations

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// PlatformType is used to tailor an instance that is being created.
type PlatformType interface {
	getPlatformName() string
}

// replaceTag replaces the specified tag values. If a tag is not found no changes are performed.
func replaceTag(input *ec2.RunInstancesInput, tags ...*ec2.Tag) {
	for _, tag := range input.TagSpecifications[0].Tags {
		for _, newTag := range tags {
			if *tag.Key == *newTag.Key {
				*tag.Value = *newTag.Value
				break
			}
		}
	}
}
