// Package connection owns outbound Integration connection events.
package connection

import (
    "context"

    mailbox "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
)

// ChangedProducer publishes a provider notification hint, not a durable outbox
// event. Gmail Pub/Sub retries its request if broker publication fails.
type ChangedProducer interface {
    PublishChanged(ctx context.Context, payload mailbox.Payload) error
}
