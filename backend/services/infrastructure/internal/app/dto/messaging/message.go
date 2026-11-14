package messaging

import coreevents "grovehub/libs/core/events"
import coremessaging "grovehub/libs/core/events/messaging"

// Internal DTOs reuse the shared transport-neutral event contracts without
// duplicating payload types or leaking internal imports through public clients.
type Message = coremessaging.Message
type MessageHandler = coremessaging.Handler
type Event = coreevents.Message
