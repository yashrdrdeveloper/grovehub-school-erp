package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidMessage = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_INVALID",
		"Mail provider message is invalid",
	)
	MessageNotFound = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_NOT_FOUND",
		"Mail provider message was not found",
	)
	MessageCreateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_CREATE_FAILED",
		"Unable to create mail provider message",
	)
	MessageReadFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_READ_FAILED",
		"Unable to read mail provider message",
	)
	MessageUpdateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_UPDATE_FAILED",
		"Unable to update mail provider message",
	)
	MessageSubmissionFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_SUBMISSION_FAILED",
		"Unable to submit mail provider message",
	)
)
