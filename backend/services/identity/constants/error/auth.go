package errorcode

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	InvalidRegistration = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_INVALID",
		"Registration request is invalid",
		http.StatusBadRequest,
	)
	InvalidRegistrationIdempotencyKey = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_IDEMPOTENCY_KEY_INVALID",
		"Registration idempotency key is invalid",
		http.StatusBadRequest,
	)
	RegistrationIdempotencyConflict = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_IDEMPOTENCY_CONFLICT",
		"Idempotency key was already used for a different registration request",
		http.StatusConflict,
	)
	RegistrationAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_ALREADY_EXISTS",
		"An identity account already exists for this registration",
		http.StatusConflict,
	)
	RegistrationFailed = coreerrors.NewDefinition(
		"IDENTITY_REGISTRATION_FAILED",
		"Unable to complete registration",
		http.StatusInternalServerError,
	)
)
