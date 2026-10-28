package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidMessage = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_INVALID",
		"Mail provider message is invalid",
	)
	InvalidMessageIdempotencyKey = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_IDEMPOTENCY_KEY_INVALID",
		"Mail provider message idempotency key is invalid",
	)
	MessageIdempotencyConflict = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_IDEMPOTENCY_CONFLICT",
		"Idempotency key was already used for a different message request",
	)
	MessageSubmissionInProgress = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_SUBMISSION_IN_PROGRESS",
		"Mail provider message submission is already in progress",
	)
	MessageRetryNotAllowed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_RETRY_NOT_ALLOWED",
		"Only failed mail provider messages can be retried",
	)
	MessageConcurrentModification = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_CONCURRENT_MODIFICATION",
		"Mail provider message was modified by another request",
	)
	InvalidMessageSearch = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_SEARCH_INVALID",
		"Mail provider message search is invalid",
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
	MessageSearchFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_MESSAGE_SEARCH_FAILED",
		"Unable to search mail provider messages",
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
