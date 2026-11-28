package constants

import messagingErrors "grovehub/services/infrastructure/internal/app/constants/error"

var (
	MessagingProviderUnavailable = messagingErrors.MessagingProviderUnavailable
	MessagingPublishFailed = messagingErrors.MessagingPublishFailed
	MessagingConsumeFailed = messagingErrors.MessagingConsumeFailed
	ErrInvalidMessage = messagingErrors.ErrInvalidMessage
	ErrInvalidStore = messagingErrors.ErrInvalidStore
	ErrReplayUnavailable = messagingErrors.ErrReplayUnavailable
)
