package messaging

import coremessaging "grovehub/libs/core/events/messaging"

// Broker-delivered message DTOs remain independent from persistence entities.
type Message = coremessaging.Message
type MessageHandler = coremessaging.Handler
