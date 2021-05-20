package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Config is used to configure a ConfigProvider.
type Config struct {
	// Region represents the region where the AWS client will connect to.
	// Specifying the wrong region may cause problems when launching simulations.
	Region string
}

// Must forces GetConfigProvider to return the ConfigProvider.
// If there is an error when calling GetConfigProvider, Must will panic.
// Example:
// 			c := Must(GetConfigProvider(Config{}))
func Must(c client.ConfigProvider, err error) client.ConfigProvider {
	if err != nil {
		panic(err)
	}
	return c
}

// GetConfigProvider returns a new Amazon Web Services session.
func GetConfigProvider(config Config) (client.ConfigProvider, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	})
}
