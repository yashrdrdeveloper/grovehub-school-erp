package messaging

import coreevents "grovehub/libs/core/events"

// Event is a transport-neutral, immutable pending event DTO.
type Event = coreevents.Message
