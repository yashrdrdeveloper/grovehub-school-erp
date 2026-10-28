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
	AccountRegistrationIdempotencyConflict = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_REGISTRATION_IDEMPOTENCY_CONFLICT",
		"Registration idempotency key is already associated with a different account request",
	)
	AccountNotFound = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_NOT_FOUND",
		"Identity account was not found",
	)
	AccountConcurrentModification = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_CONCURRENT_MODIFICATION",
		"Identity account was modified by another request",
	)
	AccountCreateFailed = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_CREATE_FAILED",
		"Unable to create identity account",
	)
	AccountReadFailed = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_READ_FAILED",
		"Unable to read identity account",
	)
	AccountUpdateFailed = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_UPDATE_FAILED",
		"Unable to update identity account",
	)
	AccountDeleteFailed = coreerrors.NewDefinition(
		"IDENTITY_ACCOUNT_DELETE_FAILED",
		"Unable to permanently delete identity account",
	)
)
