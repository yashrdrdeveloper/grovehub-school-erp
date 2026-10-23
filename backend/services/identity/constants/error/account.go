package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidBootstrapRequest = coreerrors.NewDefinition(
		"IDENTITY_BOOTSTRAP_INVALID_REQUEST",
		"Identity bootstrap request is invalid",
	)
	BootstrapAlreadyCompleted = coreerrors.NewDefinition(
		"IDENTITY_BOOTSTRAP_ALREADY_COMPLETED",
		"Identity bootstrap has already been completed",
	)
	BootstrapCreateFailed = coreerrors.NewDefinition(
		"IDENTITY_BOOTSTRAP_CREATE_FAILED",
		"Unable to create initial identity account",
	)
)
