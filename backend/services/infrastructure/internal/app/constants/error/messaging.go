package errorcode

import (
    "errors"

    coreerrors "grovehub/libs/core/errors"
)

var (
    MessagingProviderUnavailable = coreerrors.NewDefinition("INFRASTRUCTURE_MESSAGING_PROVIDER_UNAVAILABLE", "Messaging provider is unavailable")
    MessagingPublishFailed = coreerrors.NewDefinition("INFRASTRUCTURE_MESSAGING_PUBLISH_FAILED", "Unable to publish message")
    MessagingConsumeFailed = coreerrors.NewDefinition("INFRASTRUCTURE_MESSAGING_CONSUME_FAILED", "Unable to consume messages")
    ErrInvalidMessage = errors.New("messaging event requires event ID, topic and payload")
    ErrInvalidStore = errors.New("messaging repository requires a database")
)
