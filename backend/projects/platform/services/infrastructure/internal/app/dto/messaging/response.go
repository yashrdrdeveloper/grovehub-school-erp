package messaging

import inframessaging "grovehub/libs/infrastructure/messaging"

// Broker-delivered message DTOs remain independent from persistence entities.
type Message = inframessaging.Message
type MessageHandler = inframessaging.Handler
