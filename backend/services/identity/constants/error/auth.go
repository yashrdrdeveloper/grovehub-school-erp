package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidRegistration = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_INVALID",
		"Registration request is invalid",
	)
	RegistrationAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_ALREADY_EXISTS",
		"An identity account already exists for this registration",
	)
	RegistrationFailed = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_FAILED",
		"Unable to complete registration",
	)
)
