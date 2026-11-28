package messaging

import inframessaging "grovehub/services/infrastructure/client/messaging"

// Broker-delivered message DTOs remain independent from persistence entities.
type Message = inframessaging.Message
type MessageHandler = inframessaging.Handler
