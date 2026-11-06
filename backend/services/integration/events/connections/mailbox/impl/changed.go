package impl

import (
	"context"
	"encoding/json"

	infrastructureclient "grovehub/services/infrastructure/client"
	mailboxevents "grovehub/services/integration/events/connections/mailbox"
	eventhelpers "grovehub/services/integration/events/connections/mailbox/helpers"
)

// PublishChanged owns event assembly/serialization. The caller controls
// acknowledgment: a failed publish must result in a Pub/Sub retry.
func PublishChanged(ctx context.Context, messaging infrastructureclient.MessagingClient, payload mailboxevents.Payload) error {
	event := eventhelpers.BuildChanged(ctx, payload)
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return messaging.Publish(ctx, mailboxevents.ChangedTopic, encoded, eventhelpers.ChangedMetadata(event))
}
