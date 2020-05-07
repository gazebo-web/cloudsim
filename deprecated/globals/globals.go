package globals

import (
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	igntran "gitlab.com/ignitionrobotics/web/cloudsim/ign-transport"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/users"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/go-playground/form"
	"gopkg.in/go-playground/validator.v9"
	"k8s.io/client-go/kubernetes"
)

// TODO: remove as much as possible from globals

/////////////////////////////////////////////////
/// Define global constants here

/////////////////////////////////////////////////
/// Define global variables here

// Deprecated: Server encapsulates database, router, and auth0
var Server *ign.Server

// Deprecated: APIVersion is route api version.
// See also routes and routers
// \todo: Add support for multiple versions.
var APIVersion = "1.0"

// Deprecated: Validate references the global structs validator.
// See https://github.com/go-playground/validator.
// We use a single instance of validator, as it caches struct info
var Validate *validator.Validate

// Deprecated: FormDecoder holds a reference to the global Form Decoder.
// See https://github.com/go-playground/form.
// We use a single instance of Decoder, as it caches struct info
var FormDecoder *form.Decoder

// Deprecated: IgnTransport holds a reference to a ign_transport node.
var IgnTransport *igntran.GoIgnTransportNode

// Deprecated: IgnTransportTopic is the name of the topic to publish to (for testing purposes)
var IgnTransportTopic string

// Deprecated: DefaultEmailRecipients is the default recipient when sending emails.
// It is set using IGN_DEFAULT_EMAIL_RECIPIENT env var.
var DefaultEmailRecipients []string

// Deprecated: DefaultEmailSender is the default sender to use when sending emails.
// It is set using IGN_DEFAULT_EMAIL_SENDER env var.
var DefaultEmailSender string

// Deprecated: DisableSummaryEmails defines if cloudsim should send summary emails
// TODO This should probably be placed in the service configuration
var DisableSummaryEmails = false

// Deprecated: DisableScoreGeneration defines if cloud should generate score for simulations
// TODO This should probably be placed in the service configuration
var DisableScoreGeneration = false

// Deprecated: Service holds a reference to the Service. A proxy to ign-fuel's Users library
// Dev note: code should not use this from globals. Instead configure your logic with arguments
// in the constructors.
var UserAccessor useracc.UserAccessor

// Deprecated: Permissions manages permissions for users, roles and resources.
var Permissions *permissions.Permissions

// Deprecated: KClientset holds a reference to the kubernetes clientset.
// Dev note: code should not use this from globals. Instead configure your logic with arguments
// in the constructors. This is here to use from tests.
var KClientset kubernetes.Interface

// Deprecated: S3Svc holds a reference to the AWS S3 client.
// Dev note: code should not use this from globals. Instead configure your logic with arguments
// in the constructors. This is here to use from tests.
var S3Svc s3iface.S3API

// Deprecated: EC2Svc holds a reference to the AWS EC2 client.
// Dev note: code should not use this from globals. Instead configure your logic with arguments
// in the constructors. This is here to use from tests.
var EC2Svc ec2iface.EC2API
