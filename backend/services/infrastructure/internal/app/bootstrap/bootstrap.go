package bootstrap

import client "grovehub/services/infrastructure/client/messaging"

// Runtime owns provider and dispatcher lifecycle, not business services.
type Runtime interface {
    MessagingClient() client.MessagingClient
    EventStore() client.EventStore
    Close() error
}
