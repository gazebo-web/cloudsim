package fake

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/ses"
	cloud "github.com/gazebo-web/cloudsim/v4/pkg/cloud/aws"
	"github.com/gazebo-web/cloudsim/v4/pkg/defaults"
	email "github.com/gazebo-web/cloudsim/v4/pkg/email/implementations/ses"
	"github.com/gazebo-web/cloudsim/v4/pkg/machines/implementations/ec2"
	"github.com/gazebo-web/cloudsim/v4/pkg/mock"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/implementations/kubernetes"
	"github.com/gazebo-web/cloudsim/v4/pkg/platform"
	fakeSecrets "github.com/gazebo-web/cloudsim/v4/pkg/secrets/implementations/fake"
	"github.com/gazebo-web/cloudsim/v4/pkg/storage/implementations/s3"
	"github.com/gazebo-web/cloudsim/v4/pkg/store/implementations/store"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"net/http/httptest"
)

// NewInput contains input fields for the NewFakePlatform function.
type NewInput struct {
	Name    string
	Logger  gz.Logger
	Session client.ConfigProvider
	platform.Components
}

// SetDefaults sets default values
func (input *NewInput) SetDefaults() error {
	if input.Name == "" {
		input.Name = "fake"
	}

	if input.Logger == nil {
		input.Logger = gz.NewLoggerNoRollbar("fake", gz.VerbosityWarning)
	}

	if input.Session == nil {
		var err error
		if input.Session, err = session.NewSession(); err != nil {
			return err
		}
	}

	// Components
	if input.Machines == nil {
		var err error

		newInput := &ec2.NewInput{
			API:            mock.NewEC2(),
			CostCalculator: cloud.NewCostCalculatorEC2(pricing.New(input.Session)),
			Logger:         input.Logger,
			Zones: []ec2.Zone{
				{
					Zone:     "fake",
					SubnetID: "subnet-fake",
				},
			},
		}
		if input.Machines, err = ec2.NewMachines(newInput); err != nil {
			return err
		}
	}

	if input.Storage == nil {
		s3Backend := s3mem.New()
		s3Fake := gofakes3.New(s3Backend)
		s3Server := httptest.NewServer(s3Fake.Server())

		seSessionConfig := &aws.Config{
			Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
			Endpoint:         aws.String(s3Server.URL),
			Region:           aws.String("us-east-1"),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}

		s3Session, err := session.NewSession(seSessionConfig)
		if err != nil {
			return err
		}

		s3API := s3.NewAPI(s3Session)

		input.Storage = s3.NewStorage(s3API, input.Logger)
	}

	if input.Cluster == nil {
		input.Cluster, _ = kubernetes.NewFakeKubernetes(input.Logger)
	}

	if input.Store == nil {
		var err error
		input.Store, err = store.NewStoreFromEnvVars()
		if err != nil {
			return err
		}
	}

	if input.Secrets == nil {
		input.Secrets = fakeSecrets.NewFakeSecrets()
	}

	if input.EmailSender == nil {
		input.EmailSender = email.NewEmailSender(ses.New(input.Session), input.Logger)
	}

	return nil
}

// NewFakePlatform creates and returns a platform with fake components.
// If `input` or any of its fields are `nil`, default values will be used.
func NewFakePlatform(input *NewInput) (platform.Platform, error) {
	// Initialize an empty input if it was not provided
	if input == nil {
		input = &NewInput{}
	}

	if err := defaults.SetValues(input); err != nil {
		return nil, err
	}

	return platform.NewPlatform(input.Name, input.Components)
}
