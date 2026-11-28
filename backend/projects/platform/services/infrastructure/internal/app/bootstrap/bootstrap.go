package bootstrap

import client "grovehub/services/infrastructure/client/messaging"

// Runtime exposes Infrastructure's transport-neutral clients and owns its worker.
type Runtime interface {
    MessagingClient() client.Client
    EventStore() client.EventStore
    Close() error
}

type Bootstrap interface {
    Build() (Runtime, error)
}
