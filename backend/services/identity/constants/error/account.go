package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidAccount = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_INVALID",
		"Identity account is invalid",
	)
	AccountAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_ALREADY_EXISTS",
		"Identity account already exists for this person",
	)
	AccountEmailAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_EMAIL_ALREADY_EXISTS",
		"Identity account email is already in use",
	)
	AccountPhoneAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_PHONE_ALREADY_EXISTS",
		"Identity account phone is already in use",
	)
	AccountCreateFailed = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_CREATE_FAILED",
		"Unable to create identity account",
	)
)
