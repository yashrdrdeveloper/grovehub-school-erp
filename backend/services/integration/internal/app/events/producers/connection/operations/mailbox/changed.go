package mailbox

import (
    coreevents "grovehub/libs/core/events"
    flows "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox/flows"
)

// Payload describes the provider notification only.
// Communication owns message state and the authoritative sync cursor.
type Payload = flows.Payload


// Changed uses the shared generic event envelope, never a custom event wrapper.
// @Event topic=integration_mailbox_changed type=integration.mailbox.changed.v1 typeConst=ChangedEventType
type Changed = coreevents.Envelope[Payload]

