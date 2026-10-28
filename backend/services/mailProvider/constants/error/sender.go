package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidSender = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_INVALID",
		"Mail provider sender is invalid",
	)
	SenderAlreadyExists = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_ALREADY_EXISTS",
		"Mail provider sender already exists",
	)
	SenderNotFound = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_NOT_FOUND",
		"Mail provider sender was not found",
	)
	SenderConcurrentModification = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_CONCURRENT_MODIFICATION",
		"Mail provider sender was modified by another request",
	)
	SenderCreateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_CREATE_FAILED",
		"Unable to create mail provider sender",
	)
	SenderReadFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_READ_FAILED",
		"Unable to read mail provider sender",
	)
	SenderUpdateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_SENDER_UPDATE_FAILED",
		"Unable to update mail provider sender",
	)
)
