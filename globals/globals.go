package globals

import (
	"github.com/go-playground/form"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
)

// TODO: remove as much as possible from globals

/////////////////////////////////////////////////
/// Define global constants here

/////////////////////////////////////////////////
/// Define global variables here

// Server encapsulates database, router, and auth0
var Server *ign.Server

// APIVersion is route api version.
// See also routes and routers
// \todo: Add support for multiple versions.
var APIVersion = "1.0"

// Validate references the global structs validator.
// See https://github.com/go-playground/validator.
// We use a single instance of validator, as it caches struct info
var Validate *validator.Validate

// FormDecoder holds a reference to the global Form Decoder.
// See https://github.com/go-playground/form.
// We use a single instance of Decoder, as it caches struct info
var FormDecoder *form.Decoder

// DefaultEmailRecipients is the default recipient when sending emails.
// It is set using IGN_DEFAULT_EMAIL_RECIPIENT env var.
var DefaultEmailRecipients []string

// DefaultEmailSender is the default sender to use when sending emails.
// It is set using IGN_DEFAULT_EMAIL_SENDER env var.
var DefaultEmailSender string

// DisableSummaryEmails defines if cloudsim should send summary emails
// TODO This should probably be placed in the service configuration
var DisableSummaryEmails = false

// DisableScoreGeneration defines if cloud should generate score for simulations
// TODO This should probably be placed in the service configuration
var DisableScoreGeneration = false

// UserAccessor holds a reference to the UserAccessor. A proxy to ign-fuel's Users library
// Dev note: code should not use this from globals. Instead configure your logic with arguments
// in the constructors.
var UserAccessor useracc.Service

// Permissions manages permissions for users, roles and resources.
var Permissions *permissions.Permissions

// TransportTestMock holds a reference to the mock for the transport layer.
var TransportTestMock *ignws.PubSubTransporterMock
