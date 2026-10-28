package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidRegistration = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_INVALID",
		"Registration request is invalid",
	)
	InvalidRegistrationIdempotencyKey = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_IDEMPOTENCY_KEY_INVALID",
		"Registration idempotency key is invalid",
	)
	RegistrationIdempotencyConflict = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_IDEMPOTENCY_CONFLICT",
		"Idempotency key was already used for a different registration request",
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
