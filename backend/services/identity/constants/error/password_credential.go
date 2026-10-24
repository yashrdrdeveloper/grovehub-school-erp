package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidPassword = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_INVALID",
		"Password is invalid",
	)
	PasswordCredentialAlreadyExists = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_ALREADY_EXISTS",
		"Password credential already exists for this account",
	)
	PasswordCredentialNotFound = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_NOT_FOUND",
		"Password credential was not found",
	)
	PasswordCredentialCreateFailed = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_CREATE_FAILED",
		"Unable to create password credential",
	)
	PasswordCredentialReadFailed = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_READ_FAILED",
		"Unable to read password credential",
	)
	PasswordCredentialUpdateFailed = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_UPDATE_FAILED",
		"Unable to update password credential",
	)
	PasswordCredentialVerifyFailed = coreerrors.NewDefinition(
		"IDENTITY_PASSWORD_CREDENTIAL_VERIFY_FAILED",
		"Unable to verify password credential",
	)
)
